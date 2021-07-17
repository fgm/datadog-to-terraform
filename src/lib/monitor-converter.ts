import { assignmentString, block, convertFromDefinition, map } from "./utils";

const MONITOR = {
  message: (v: any) => assignmentString("message", v),
  name: (v: any) => assignmentString("name", v),
  query: (v: any) => assignmentString("query", v),
  type: (v: any) => assignmentString("type", v),
  options: (v: any) => map(v, (k1, v1) => convertFromDefinition(OPTIONS, k1, v1)),
  id: (_: any) => "",
  tags: (v: any) => assignmentString("tags", v),
  priority: (v: any) => assignmentString("priority", v),
};

const OPTIONS = {
  enable_logs_sample: (v: any) => assignmentString("enable_logs_sample", v),
  escalation_message: (v: any) => assignmentString("escalation_message", v),
  evaluation_delay: (v: any) => assignmentString("evaluation_delay", v),
  force_delete: (v: any) => assignmentString("force_delete", v),
  include_tags: (v: any) => assignmentString("include_tags", v),
  locked: (v: any) => assignmentString("locked", v),
  new_host_delay: (v: any) => assignmentString("new_host_delay", v),
  no_data_timeframe: (v: any) => assignmentString("no_data_timeframe", v),
  notify_audit: (v: any) => assignmentString("notify_audit", v),
  notify_no_data: (v: any) => assignmentString("notify_no_data", v),
  renotify_interval: (v: any) => assignmentString("renotify_interval", v),
  require_full_window: (v: any) => assignmentString("require_full_window", v),
  restricted_roles: (v: any) => assignmentString("restricted_roles", v),
  threshold_windows: (v: any) => block("monitor_threshold_windows", v, assignmentString),
  thresholds: (v: any) => block("monitor_thresholds", v, assignmentString),
  timeout_h: (v: any) => assignmentString("timeout_h", v),
  validate: (v: any) => assignmentString("timeout_h", v),
  groupby_simple_monitor: (v: any) => assignmentString("groupby_simple_monitor", v),
  silenced: (_: any) => "", // 2.23.0 deprecated
};

export function generateTerraformCode(resourceName: string, monitorData: object) {
  let result = "";
  Object.entries(monitorData).forEach(([k, v]) => {
    result += convertFromDefinition(MONITOR, k, v);
  });
  return `resource "datadog_monitor" "${resourceName}" {${result}\n}`;
}
