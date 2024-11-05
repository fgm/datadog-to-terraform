package main

import (
	"fmt"
	"strings"
)

// MONITOR and OPTIONS definitions
var MONITOR = map[string]stringFunc{
	"message": stringGen("message"),
	"name":    stringGen("name"),
	"query":   stringGen("query"),
	"type":    stringGen("type"),
	"options": func(v any) string {
		return mapContents(v.(map[string]any), func(k1 string, v1 any) string { return convertFromDefinition(OPTIONS, k1, v1) })
	},
	"id":               blankGen,
	"tags":             stringGen("tags"),
	"priority":         stringGen("priority"),
	"restricted_roles": stringGen("restricted_roles"),
}

var OPTIONS = map[string]stringFunc{
	"enable_logs_sample":  stringGen("enable_logs_sample"),
	"escalation_message":  stringGen("escalation_message"),
	"evaluation_delay":    stringGen("evaluation_delay"),
	"force_delete":        stringGen("force_delete"),
	"include_tags":        stringGen("include_tags"),
	"locked":              stringGen("locked"),
	"new_host_delay":      stringGen("new_host_delay"),
	"no_data_timeframe":   stringGen("no_data_timeframe"),
	"notify_audit":        stringGen("notify_audit"),
	"notify_no_data":      stringGen("notify_no_data"),
	"renotify_interval":   stringGen("renotify_interval"),
	"require_full_window": stringGen("require_full_window"),
	"restricted_roles":    stringGen("restricted_roles"),
	"threshold_windows": func(v any) string {
		return block("monitor_threshold_windows", v.(map[string]any), assignmentString)
	},
	"thresholds": func(v any) string {
		return block("monitor_thresholds", v.(map[string]any), assignmentString)
	},
	"timeout_h":              stringGen("timeout_h"),
	"validate":               stringGen("timeout_h"),
	"groupby_simple_monitor": stringGen("groupby_simple_monitor"),
	"silenced":               blankGen, // Deprecated
	"new_group_delay":        stringGen("new_group_delay"),
	"renotify_statuses":      stringGen("renotify_statuses"),
	"on_missing_data":        stringGen("on_missing_data"),
}

func generateMonitorTerraformCode(resourceName string, monitorData map[string]any) string {
	var result strings.Builder
	for k, v := range monitorData {
		result.WriteString(convertFromDefinition(MONITOR, k, v))
	}
	return fmt.Sprintf("resource \"datadog_monitor\" \"%s\" {%s\n}", resourceName, result.String())
}
