package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fgm/datadog-to-terraform/convert"
)

func readInput(path string) (convert.JSONData, convert.DataDogDocumentType, error) {
	var bs []byte
	var err error

	if path == "-" {
		bs, err = io.ReadAll(os.Stdin)
		if err != nil {
			return nil, convert.InvalidType, fmt.Errorf("reading stdin: %w\n", err)
		}
	} else {
		bs, err = os.ReadFile(path)
		if err != nil {
			return nil, convert.InvalidType, fmt.Errorf("reading %s: %w\n", path, err)
		}
	}

	jd := make(convert.JSONData)
	if err = json.Unmarshal(bs, &jd); err != nil {
		log.Fatalf("decoding JSON from input: %v\n", err)
	}
	if _, ok := jd["name"]; ok {
		return jd, convert.MonitorType, nil
	}
	return jd, convert.DashboardType, nil
}

func emitOutput(path string, tf string) error {
	var err error

	if path == "-" {
		_, err = os.Stdout.WriteString(tf)
		if err != nil {
			return fmt.Errorf("writing to stdout: %w", err)
		}
	} else {
		err = os.WriteFile(path, []byte(tf), 0666)
		if err != nil {
			return fmt.Errorf("writing to %s: %w", path, err)
		}
	}
	return nil
}

func main() {
	inFile := flag.String("if", "-", `Input file. "-" means standard input.`)
	outFile := flag.String("of", "-", `Output file. "-" means standard output.`)
	flag.Parse()

	var dt convert.DataDogDocumentType
	var err error
	var jd convert.JSONData
	var tf string

	jd, dt, err = readInput(*inFile)
	if err != nil {
		log.Fatalln(err)
	}

	switch dt {
	case convert.DashboardType:
		tf, err = convert.Dashboard(jd)
	case convert.MonitorType:
		tf, err = convert.Monitor(jd)
	default:
		tf, err = "", errors.New("unrecognized JSON document type")
	}
	if err != nil {
		log.Fatalf("Conversion failed: %v\n", err)
	}

	err = emitOutput(*outFile, tf)
	if err != nil {
		log.Fatalln(err)
	}
}
