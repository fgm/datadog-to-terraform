package convert

import (
	"fmt"
	"strings"
)

func boolVal(v interface{}) string {
	b := v.(bool)
	if b {
		return "true"
	}
	return "false"
}

func droppedVal(_ interface{}) string {
	return `"dropped"`
}

func intListVal(v interface{}) string {
	ss := v.([]string)
	build := []string{"{"}
	for _, s := range ss {
		build = append(build, indent+intVal(s)+",")
	}
	build = append(build, "}")
	return strings.Join(build, "\n")
}

func intVal(v interface{}) string {
	return fmt.Sprintf("%d", v)
}

func stringListVal(v interface{}) string {
	ss := v.([]string)
	build := []string{"{"}
	for _, s := range ss {
		build = append(build, indent+stringVal(s)+",")
	}
	build = append(build, "}")
	return strings.Join(build, "\n")
}

func stringVal(v interface{}) string {
	return fmt.Sprintf(`"%s"`, v)
}
