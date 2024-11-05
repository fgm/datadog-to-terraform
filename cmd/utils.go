package main

import (
	"fmt"
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
	keys := make([]string, 0, len(contents))
	for k := range contents {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		result.WriteString(converter(k, contents[k]))
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
func convertFromDefinition(definitionSet map[string]stringFunc, name string, v any) (string, error) {
	if converter, exists := definitionSet[name]; exists {
		return converter(v), nil
	}
	return "", fmt.Errorf("can't convert key '%s' with value %#v", name, v)
}

// jmapsFromAny extract a jmaps from an "any" value which is actually a jmap
// or a slice in which elements are also "any" values with a dynamic jmap value.
func jmapsFromAny(v any) (jmaps, error) {
	slice, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("items expected as []any but got %T: %#v", v, v)
	}
	items := make(jmaps, len(slice))
	for i, item := range slice {
		items[i], ok = item.(jmap)
		if !ok {
			return nil, fmt.Errorf("item [%d] expected as jmap but got %T: %#v", i, item, item)
		}
	}
	return items, nil
}

func convertSlice[T any](v []T) string {
	var result strings.Builder
	result.WriteString("[")
	max := len(v) - 1
	for i, elem := range v {
		result.WriteString(literalString(elem))
		if i != max {
			result.WriteString(",")
		}
	}
	result.WriteString("]")
	return result.String()
}

// Converts a string or list of strings value to a literal string representation
func literalString(value any) string {
	switch v := value.(type) {
	case string:
		if strings.Contains(value.(string), "\n") {
			return fmt.Sprintf("<<EOF\n%s\nEOF", v)
		}
		return fmt.Sprintf("\"%s\"", v)
	case []any:
		return convertSlice(v)
	case []string:
		return convertSlice(v)
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

func must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
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
