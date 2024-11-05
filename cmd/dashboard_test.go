package main

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/sebdah/goldie/v2"
)

//go:embed testdata/dashboard-good.json
var dashGoodJSON []byte

//go:embed testdata/dashboard-bad.json
var dashBadJSON []byte

//go:embed testdata/screenboard.json
var screenboardJSON []byte

//go:embed testdata/timeboard.json
var timeboardJSON []byte

func Test_generateDashboardTerraformCode(t *testing.T) {
	for _, test := range [...]struct {
		name    string
		input   []byte
		errorRx string
	}{
		{"dashboard-good", dashGoodJSON, ""},
		{"dashboard-bad", dashBadJSON, "^can't convert key.*with value \".*\"$"},
		{"screenboard", screenboardJSON, ""},
		{"timeboard", timeboardJSON, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			var j = make(jmap, 0)
			if err := json.Unmarshal(test.input, &j); err != nil {
				t.Fatal(err)
			}

			g := goldie.New(t)
			actual, err := generateDashboardTerraformCode(test.name, j)
			switch {
			case test.errorRx == "" && err != nil:
				t.Fatalf("unexpected error %v", err)
			case test.errorRx != "" && err == nil:
				t.Fatalf("unexpected success")
			case test.errorRx != "" && err != nil:
				rx := regexp.MustCompile(test.errorRx)
				if !rx.MatchString(err.Error()) {
					t.Errorf("got error %s but expected to match %s", err, rx.String())
				}
			}
			g.Assert(t, test.name, []byte(actual))
		})
	}
}
