package convert

import (
	"fmt"
)

type dashboard struct {
	Name string
	TFBlock
}

func (d *dashboard) Render(level int, name string) (string, error) {
	name = fmt.Sprintf(`datadog_dashboard "%s" {`, d.Name)
	return d.TFBlock.Render(0, name)
}

func Dashboard(jd JSONData, name string) (string, error) {
	db := TFBlock{
		JSONTag:    "dashboard",
		RootTag:    "dashboard",
		props:      make(map[string]Renderable, len(jd)), // Most values are props, not blocks
		blockLists: make(map[string]Renderable, 3),       // widget, template_variable, template_variable_preset
	}
	if err := db.Load(jd, [2]string{db.RootTag, db.JSONTag}); err != nil {
		return "", fmt.Errorf("loading dashboard: %w", err)
	}

	return db.Render(0, fmt.Sprintf(`resource "datadog_dashboard" "%s"`, name))
}

func makeBlockVal(name string, schemaPath []string) Converter {
	return func(name string, v interface{}) (Renderable, error) {
		return stringValue(""), nil
	}
}
