package convert

func Dashboard(jd JSONData) (string, error) {
	return "", nil
}

type dashboard struct {
}

type DocElement struct {
	Converter
	Required bool
	Hidden   bool
}

func (de *DocElement) IndentedString(indent int) string {

}

func makeBlockVal(name string) Converter {
	return func(v interface{}) string {

	}
}

// Keys are the names of the keys in the JSON version.
// They usually do NOT match the TF block names, e.g "widgets" in JSON
// translates to a list of "widget" blocks with equivalent content; while
// "queries" in JSON translates to a list of "query" blocks each containing a
// single "metric_query" block with equivalent content.
var dashboardSchema = map[string]map[string]DocElement{
	"dashboard": {
		// Required
		"layout_type": {stringVal, true, false},
		"title":       {stringVal, true, false},
		"widget":      {makeBlockVal("widget"), true, false},

		// Normal
		"dashboard_lists":  {intListVal, false, false},
		"description":      {stringVal, false, false},
		"is_read_only":     {boolVal, false, false},
		"notify_list":      {stringListVal, false, false},
		"reflow_type":      {stringVal, false, false}, // "auto"|"fixed". Only if layout_type is "ordered"
		"restricted_roles": {stringListVal, false, false},
		"url":              {stringVal, false, true},

		// Hidden
		"dashboard_lists_removed": {droppedVal, false, true},
		"id":                      {stringVal, false, true},

		// TODO
		// "template_variable" []Block,
		// "template_variable_preset" []Block
	},
	"wdigets":,
}
