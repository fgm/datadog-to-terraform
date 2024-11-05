package main

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
)

type (
	jmap       = map[string]any
	jmaps      = []jmap
	stringFunc func(any) string
)

var (
	blankGen = stringFunc(func(_ any) string { return "" })
)

func stringGen(name string) stringFunc {
	return func(v any) string { return assignmentString(name, v) }
}

var WIDGET_DEFINITION map[string]stringFunc

var DASHBOARD = map[string]stringFunc{
	"dashboard_lists_removed": blankGen,
	"description":             stringGen("description"),
	"id":                      blankGen,
	"is_read_only":            stringGen("is_read_only"),
	"layout_type":             stringGen("layout_type"),
	"notify_list":             stringGen("notify_list"),
	"reflow_type":             stringGen("reflow_type"),
	"restricted_roles":        func(_ any) string { return "// restricted_roles is still in beta" },
	"template_variables": func(v any) string {
		slice, ok := v.([]any)
		if !ok {
			log.Fatalf("template_variables expected as []any but got %T: %#v\n", v, v)
		}
		tvs := make(jmaps, len(slice))
		for i, tv := range slice {
			tvs[i], ok = tv.(jmap)
			if !ok {
				log.Fatalf("template_variables[%d] expected as jmaps but got %T: %#v\n", i, tv, tv)
			}
		}
		return blockList(tvs, "template_variable", assignmentString)
	},
	"template_variable_presets": func(v any) string {
		presets := must(jmapsFromAny(v))
		return blockList(presets, "template_variable_preset", func(k1 string, v1 any) string {
			return must(convertFromDefinition(TEMPLATE_VARIABLE_PRESET, k1, v1))
		})
	},
	"title": stringGen("title"),
	"url":   stringGen("url"),
	"widgets": func(v any) string {
		widgets := must(jmapsFromAny(v))
		return convertWidgets(widgets)
	},
}

var EVENT_QUERY = map[string]stringFunc{
	"aggregator":  stringGen("aggregator"),
	"compute":     func(v any) string { return block("compute", v.(jmap), assignmentString) },
	"data_source": stringGen("data_source"),
	"group_by": func(v any) string {
		groups := must(jmapsFromAny(v))
		return blockList(groups, "group_by", func(k1 string, v1 any) string {
			return must(convertFromDefinition(EVENT_QUERY_GROUP_BY, k1, v1))
		})
	},
	"indexes": stringGen("indexes"),
	"name":    stringGen("name"),
	"search":  func(v any) string { return block("search", v.(jmap), assignmentString) },
}

var EVENT_QUERY_GROUP_BY = map[string]stringFunc{
	"facet": stringGen("facet"),
	"limit": stringGen("limit"),
	"sort":  func(v any) string { return block("sort", v.(jmap), assignmentString) },
}

var FORMULA = map[string]stringFunc{
	"alias":   stringGen("alias"),
	"formula": stringGen("formula_expression "),
	"limit": func(v any) string {
		return block("limit", v.(jmap), func(k1 string, v1 any) string {
			return must(convertFromDefinition(FORMULA_LIMIT, k1, v1))
		})
	},
}

var FORMULA_LIMIT = map[string]stringFunc{
	"count": stringGen("count"),
	"order": stringGen("order"),
}

var GROUP_BY = map[string]stringFunc{
	"facet":      stringGen("facet"),
	"limit":      stringGen("limit"),
	"sort":       func(v any) string { return block("sort_query", v.(jmap), assignmentString) },
	"sort_query": func(v any) string { return block("sort_query", v.(jmap), assignmentString) },
}

var LOG_QUERY = map[string]stringFunc{
	"compute": func(v any) string {
		return block("compute_query", v.(jmap), assignmentString)
	},
	"group_by": func(v any) string {
		groups := must(jmapsFromAny(v))
		return blockList(groups, "group_by", func(k1 string, v1 any) string {
			return must(convertFromDefinition(GROUP_BY, k1, v1))
		})
	},
	"index": stringGen("index"),
	"multi_compute": func(v any) string {
		comps := must(jmapsFromAny(v))
		return blockList(comps, "multi_compute", assignmentString)
	},
	"search": func(v any) string {
		return assignmentString("search_query", v.(jmap)["query"])
	},
	"search_query": stringGen("search_query"),
}

var QUERY = map[string]stringFunc{
	"name":        stringGen("name"),
	"data_source": stringGen("data_source"),
	"query":       stringGen("query"),
}

var REQUEST = map[string]stringFunc{
	"aggregator":        stringGen("aggregator"),
	"alias":             stringGen("alias"),
	"apm_query":         stringGen("apm_query"),
	"apm_stats_query":   stringGen("apm_stats_query"),
	"cell_display_mode": stringGen("cell_display_mode"),
	"change_type":       stringGen("change_type"),
	"compare_to":        stringGen("compare_to"),
	"conditional_formats": func(v any) string {
		formats := must(jmapsFromAny(v))
		return blockList(formats, "conditional_formats", assignmentString)
	},
	"display_type": stringGen("display_type"),
	"fill":         func(v any) string { return block("fill", v.(jmap), assignmentString) },
	"formulas": func(v any) string {
		fs := must(jmapsFromAny(v))
		return blockList(fs, "formula", func(k1 string, v1 any) string {
			return must(convertFromDefinition(FORMULA, k1, v1))
		})
	},
	"increase_good": stringGen("increase_good"),
	"limit":         stringGen("limit"),
	"log_query": func(v any) string {
		return block("log_query", v.(jmap), func(k1 string, v1 any) string {
			return must(convertFromDefinition(LOG_QUERY, k1, v1))
		})
	},
	"metadata": func(v any) string {
		meta := must(jmapsFromAny(v))
		return blockList(meta, "metadata", assignmentString)
	},
	"network_query":  stringGen("network_query"),
	"on_right_yaxis": stringGen("on_right_yaxis"),
	"order":          stringGen("order"),
	"order_by":       stringGen("order_by"),
	"order_dir":      stringGen("order_dir"),
	"process_query":  stringGen("process_query"),
	"q":              stringGen("q"),
	"queries": func(v any) string {
		queries := must(jmapsFromAny(v))
		return queryBlockList(queries, assignmentString)
	},
	"response_format": stringGen("response_format"),
	"rum_query":       stringGen("rum_query"),
	"security_query":  stringGen("security_query"),
	"show_present":    stringGen("show_present"),
	"style": func(v any) string {
		return blockList(jmaps{v.(jmap)}, "style", assignmentString)
	},
}

var TEMPLATE_VARIABLE_PRESET = map[string]stringFunc{
	"name": stringGen("name"),
	"template_variables": func(v any) string {
		vars := must(jmapsFromAny(v))
		return blockList(vars, "template_variable", assignmentString)
	},
}

var WIDGET = map[string]stringFunc{
	"definition": func(v any) string { return widgetDefinition(v.(jmap)) },
	"id":         blankGen,
	"layout": func(v any) string {
		return block("widget_layout", v.(jmap), assignmentString)
	},
}

func init() {
	WIDGET_DEFINITION = map[string]stringFunc{
		"alert_id":         stringGen("alert_id"),
		"autoscale":        stringGen("autoscale"),
		"background_color": stringGen("background_color"),
		"check":            stringGen("check"),
		"color":            stringGen("color"),
		"color_by_groups":  stringGen("color_by_groups"),
		"color_preference": stringGen("color_preference"),
		"columns":          stringGen("columns"),
		"content":          stringGen("content"),
		"count":            blankGen,
		"custom_links": func(v any) string {
			links := must(jmapsFromAny(v))
			return blockList(links, "custom_link", assignmentString)
		},
		"custom_unit":    stringGen("custom_unit"),
		"display_format": stringGen("display_format"),
		"env":            stringGen("env"),
		"event":          func(v any) string { return block("event", v.(jmap), assignmentString) },
		"events": func(v any) string {
			events := must(jmapsFromAny(v))
			return blockList(events, "event", assignmentString)
		},
		"event_size":         stringGen("event_size"),
		"filters":            stringGen("filters"),
		"font_size":          stringGen("font_size"),
		"global_time_target": stringGen("global_time_target"),
		"group":              stringGen("group"),
		"group_by":           stringGen("group_by"),
		"grouping":           stringGen("grouping"),
		"has_padding":        blankGen,
		"has_search_bar":     stringGen("has_search_bar"),
		"hide_zero_counts":   stringGen("hide_zero_counts"),
		"indexes":            stringGen("indexes"),
		"layout_type":        stringGen("layout_type"),
		"legend_columns":     stringGen("legend_columns"),
		"legend_layout":      stringGen("legend_layout"),
		"legend_size":        stringGen("legend_size"),
		"live_span":          stringGen("live_span"),
		"logset":             blankGen,
		"margin":             stringGen("margin"),
		"markers": func(v any) string {
			markers := must(jmapsFromAny(v))
			return blockList(markers, "marker", assignmentString)
		},
		"message_display":     stringGen("message_display"),
		"no_group_hosts":      stringGen("no_group_hosts"),
		"no_metric_hosts":     stringGen("no_metric_hosts"),
		"node_type":           stringGen("node_type"),
		"precision":           stringGen("precision"),
		"query":               stringGen("query"),
		"requests":            convertRequests,
		"right_yaxis":         func(v any) string { return block("right_yaxis", v.(jmap), assignmentString) },
		"scope":               stringGen("scope"),
		"service":             stringGen("service"),
		"show_breakdown":      stringGen("show_breakdown"),
		"show_date_column":    stringGen("show_date_column"),
		"show_distribution":   stringGen("show_distribution"),
		"show_error_budget":   stringGen("show_error_budget"),
		"show_errors":         stringGen("show_errors"),
		"show_hits":           stringGen("show_hits"),
		"show_last_triggered": stringGen("show_last_triggered"),
		"show_latency":        stringGen("show_latency"),
		"show_legend":         stringGen("show_legend"),
		"show_message_column": stringGen("show_message_column"),
		"show_resource_list":  stringGen("show_resource_list"),
		"show_tick":           stringGen("show_tick"),
		"size_format":         stringGen("size_format"),
		"sizing":              stringGen("sizing"),
		"slo_id":              stringGen("slo_id"),
		"sort":                convertSort,
		"span_name":           stringGen("span_name"),
		"start":               blankGen,
		"style":               func(v any) string { return block("style", v.(jmap), assignmentString) },
		"summary_type":        stringGen("summary_type"),
		"tags":                stringGen("tags"),
		"tags_execution":      stringGen("tags_execution"),
		"text":                stringGen("text"),
		"text_align":          stringGen("text_align"),
		"tick_edge":           stringGen("tick_edge"),
		"tick_pos":            stringGen("tick_pos"),
		"time": func(v any) string {
			if liveSpan, ok := v.(jmap)["live_span"]; ok {
				return assignmentString("live_span", liveSpan)
			}
			return ""
		},
		"time_windows":   stringGen("time_windows"),
		"title":          stringGen("title"),
		"title_align":    stringGen("title_align"),
		"title_size":     stringGen("title_size"),
		"type":           blankGen,
		"unit":           stringGen("unit"),
		"url":            stringGen("url"),
		"vertical_align": blankGen,
		"view_mode":      stringGen("view_mode"),
		"view_type":      stringGen("view_type"),
		"viz_type":       stringGen("viz_type"),
		"widget_layout": func(v any) string {
			return block("widget_layout", v.(jmap), assignmentString)
		},
		"widgets": func(v any) string {
			return convertWidgets(must(jmapsFromAny(v)))
		},
		"xaxis": func(v any) string { return block("xaxis", v.(jmap), assignmentString) },
		"yaxis": func(v any) string { return block("yaxis", v.(jmap), assignmentString) },
	}
}

func convertEventQuery(value jmap) string {
	return block("query", value, func(_ string, _ any) string {
		return blockList(jmaps{value}, "event_query", func(k1 string, v1 any) string {
			return must(convertFromDefinition(EVENT_QUERY, k1, v1))
		})
	})
}

// convertRequests accepts either a single request as a jmap or a requests jmaps.
func convertRequests(value any) string {
	if reflect.ValueOf(value).Kind() == reflect.Slice {
		values := must(jmapsFromAny(value))
		return blockList(values, "request", func(k string, v any) string {
			return must(convertFromDefinition(REQUEST, k, v))
		})
	}
	return block("request", value.(jmap), func(k string, v any) string {
		return must(convertFromDefinition(REQUEST, k, v))
	})
}

func convertSort(v any) string {
	if sortStr, ok := v.(string); ok {
		return assignmentString("sort", sortStr)
	}
	return block("sort", v.(jmap), assignmentString)
}

func convertWidgets(value jmaps) string {
	return blockList(value, "widget", func(k1 string, v1 any) string {
		return must(convertFromDefinition(WIDGET, k1, v1))
	})
}

func widgetDefinition(contents jmap) string {
	definitionType := contents["type"].(string)
	if definitionType == "slo" {
		definitionType = "service_level_objective"
	}
	return block(fmt.Sprintf("\n%s_definition", definitionType), contents, func(k string, v any) string {
		return must(convertFromDefinition(WIDGET_DEFINITION, k, v))
	})
}

func generateDashboardTerraformCode(resourceName string, data jmap) (string, error) {
	var (
		result strings.Builder
		keys   = make([]string, 0, len(data))
	)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s, err := convertFromDefinition(DASHBOARD, k, data[k])
		if err != nil {
			return "", err
		}
		result.WriteString(s)
	}
	return fmt.Sprintf("resource \"datadog_dashboard\" \"%s\" {%s\n}", resourceName, result.String()), nil
}
