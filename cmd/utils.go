package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// Generates an assignment string for a key-value pair
func assignmentString(key string, value any) string {
	if value == nil {
		return ""
	}
	displayValue := literalString(value)
	return fmt.Sprintf("\n%s = %s", key, displayValue)
}

// Creates a block with a name and converted contents
func block(name string, contents jmap, converter func(string, any) string) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("\n%s {", name))
	for k, v := range contents {
		result.WriteString(converter(k, v))
	}
	result.WriteString("\n}")
	return result.String()
}

// Generates a list of blocks
func blockList(array jmaps, blockName string, contentConverter func(string, any) string) string {
	var result strings.Builder
	result.WriteString("\n")
	for _, elem := range array {
		result.WriteString(block(blockName, elem, contentConverter))
	}
	return result.String()
}

// Converts a key-value pair using a definition set
func convertFromDefinition(definitionSet map[string]stringFunc, name string, v any) string {
	if converter, exists := definitionSet[name]; exists {
		return converter(v)
	}
	log.Fatalf("Can't convert key '%s' with value %#v\n", name, v)
	return ""
}

func jmapsFromAny(v any) jmaps {
	slice, ok := v.([]any)
	if !ok {
		log.Fatalf("widgets expected as []any but got %T: %#v\n", v, v)
	}
	widgets := make(jmaps, len(slice))
	for i, widget := range slice {
		widgets[i], ok = widget.(jmap)
		if !ok {
			log.Fatalf("widgets[%d] expected as jmap but got %T: %#v\n", i, widget, widget)
		}
	}
	return widgets
}

// Converts a value to a literal string representation
func literalString(value any) string {
	switch v := value.(type) {
	case string:
		if strings.Contains(v, "\n") {
			return fmt.Sprintf("<<EOF\n%s\nEOF", v)
		}
		return fmt.Sprintf("\"%s\"", v)
	case []any:
		var result strings.Builder
		result.WriteString("[")
		for i, elem := range v {
			result.WriteString(literalString(elem))
			if i != len(v)-1 {
				result.WriteString(",")
			}
		}
		result.WriteString("]")
		return result.String()
	default:
		return fmt.Sprintf("%v", value)
	}
}

// Maps over contents and applies a converter function
func mapContents(contents jmap, converter func(string, any) string) string {
	var result strings.Builder
	keys := make([]string, 0, len(contents))
	for k := range contents {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		result.WriteString(converter(k, contents[k]))
	}
	return result.String()
}

// Creates a query block with a name and converted contents
func queryBlock(name string, contents jmap, converter func(string, any) string) string {
	return fmt.Sprintf("\nquery {\n\n  %s {%s\n}}", name, mapContents(contents, converter))
}

// Generates a list of query blocks
func queryBlockList(array jmaps, contentConverter func(string, any) string) string {
	var result []string
	result = append(result, "\n")
	for _, elem := range array {
		result = append(result, queryBlock("metric_query", elem, contentConverter))
	}
	return strings.Join(result, "")
}
