import {assignmentString, block, blockList, convertFromDefinition} from "./utils.js";

const DASHBOARD = {
    description: (v: object) => assignmentString("description", v),
    id: (_: any) => "",
    is_read_only: (v: object) => assignmentString("is_read_only", v),
    layout_type: (v: object) => assignmentString("layout_type", v),
    notify_list: (v: object) => assignmentString("notify_list", v),
    template_variable_presets: (v: object[]) => blockList(v, "template_variable_preset", (k1, v1) => convertFromDefinition(TEMPLATE_VARIABLE_PRESET, k1, v1)),
    template_variables: (v: object[]) => blockList(v, "template_variable", assignmentString),
    title: (v: object) => assignmentString("title", v),
    url: (v: object) => assignmentString("url", v),
    widgets: (v: object[]) => convertWidgets(v),
};

const WIDGET = {
    definition: (v: object) => widgetDefinition(v),
    id: (_: any) => "",
    layout: (v: object) => block("widget_layout", v, assignmentString),
};

const TEMPLATE_VARIABLE_PRESET = {
    name: (v: object) => assignmentString("name", v),
    template_variables: (v: object[]) => blockList(v, "template_variable", assignmentString),
};

const WIDGET_DEFINITION = {
    alert_id: (v: object) => assignmentString("alert_id", v),
    autoscale: (v: object) => assignmentString("autoscale", v),
    background_color: (v: object) => assignmentString("background_color", v),
    check: (v: object) => assignmentString("check", v),
    color: (v: object) => assignmentString("color", v),
    color_by_groups: (v: object) => assignmentString("color_by_groups", v),
    color_preference: (v: object) => assignmentString("color_preference", v),
    columns: (v: object) => assignmentString("columns", v),
    content: (v: object) => assignmentString("content", v),
    count: (_: any) => "", // 2.23.0 deprecated, see docs for widget.manage_status_definition
    custom_links: (v: object[]) => blockList(v, "custom_link", assignmentString),
    custom_unit: (v: object) => assignmentString("custom_unit", v),
    display_format: (v: object) => assignmentString("display_format", v),
    env: (v: object) => assignmentString("env", v),
    event: (v: object) => block("event", v, assignmentString),
    event_size: (v: object) => assignmentString("event_size", v),
    filters: (v: object) => assignmentString("filters", v),
    font_size: (v: object) => assignmentString("font_size", v),
    global_time_target: (_: any) => "",
    group: (v: object) => assignmentString("group", v),
    group_by: (v: object) => assignmentString("group_by", v),
    grouping: (v: object) => assignmentString("grouping", v),
    has_padding: (_: any) => "", // 2.23.0 not described in docs, occurs in widget.note_definition json
    has_search_bar: (v: object) => assignmentString("has_search_bar", v),
    hide_zero_counts: (v: object) => assignmentString("hide_zero_counts", v),
    indexes: (v: object) => assignmentString("indexes", v),
    layout_type: (v: object) => assignmentString("layout_type", v),
    legend_columns: (v: object) => assignmentString("legend_columns", v),
    legend_layout: (v: object) => assignmentString("legend_layout", v),
    legend_size: (v: object) => assignmentString("legend_size", v),
    live_span: (v: object) => assignmentString("live_span", v),
    logset: (_: any) => "", // 2.23.0 deprecated, see docs for widget.log_stream_definition
    margin: (v: object) => assignmentString("margin", v),
    markers: (v: object[]) => blockList(v, "marker", assignmentString),
    message_display: (v: object) => assignmentString("message_display", v),
    no_group_hosts: (v: object) => assignmentString("no_group_hosts", v),
    no_metric_hosts: (v: object) => assignmentString("no_metric_hosts", v),
    node_type: (v: object) => assignmentString("node_type", v),
    precision: (v: object) => assignmentString("precision", v),
    query: (v: object) => assignmentString("query", v),
    requests: (v: object) => convertRequests(v),
    right_yaxis: (v: object) => block("right_yaxis", v, assignmentString),
    scope: (v: object) => assignmentString("scope", v),
    service: (v: object) => assignmentString("service", v),
    show_breakdown: (v: object) => assignmentString("show_breakdown", v),
    show_date_column: (v: object) => assignmentString("show_date_column", v),
    show_distribution: (v: object) => assignmentString("show_distribution", v),
    show_error_budget: (v: object) => assignmentString("show_error_budget", v),
    show_errors: (v: object) => assignmentString("show_errors", v),
    show_hits: (v: object) => assignmentString("show_hits", v),
    show_last_triggered: (v: object) => assignmentString("show_last_triggered", v),
    show_latency: (v: object) => assignmentString("show_latency", v),
    show_legend: (v: object) => assignmentString("show_legend", v),
    show_message_column: (v: object) => assignmentString("show_message_column", v),
    show_resource_list: (v: object) => assignmentString("show_resource_list", v),
    show_tick: (v: object) => assignmentString("show_tick", v),
    size_format: (v: object) => assignmentString("size_format", v),
    sizing: (v: object) => assignmentString("sizing", v),
    slo_id: (v: object) => assignmentString("slo_id", v),
    sort: (v: object) => convertSort(v),
    span_name: (v: object) => assignmentString("span_name", v),
    start: (_: any) => "", // 2.23.0 deprecated, see docs for widget.manage_status_definition
    style: (v: object) => block("style", v, assignmentString),
    summary_type: (v: object) => assignmentString("summary_type", v),
    tags: (v: object) => assignmentString("tags", v),
    tags_execution: (v: object) => assignmentString("tags_execution", v),
    text: (v: object) => assignmentString("text", v),
    text_align: (v: object) => assignmentString("text_align", v),
    tick_edge: (v: object) => assignmentString("tick_edge", v),
    tick_pos: (v: object) => assignmentString("tick_pos", v),
    time: (v: any) => (!!v.live_span ? assignmentString("live_span", v.live_span) : ""),
    time_windows: (v: object) => assignmentString("time_windows", v),
    title: (v: object) => assignmentString("title", v),
    title_align: (v: object) => assignmentString("title_align", v),
    title_size: (v: object) => assignmentString("title_size", v),
    type: (_: any) => "",
    unit: (v: object) => assignmentString("unit", v),
    url: (v: object) => assignmentString("url", v),
    vertical_align: (_: any) => "", // 2.23.0 not described in docs, occurs in widget.note_definition json
    view_mode: (v: object) => assignmentString("view_mode", v),
    view_type: (v: object) => assignmentString("view_type", v),
    viz_type: (v: object) => assignmentString("viz_type", v),
    widget_layout: (v: object) => block("widget_layout", v, assignmentString),
    widgets: (v: object[]) => convertWidgets(v),
    xaxis: (v: object) => block("xaxis", v, assignmentString),
    yaxis: (v: object) => block("yaxis", v, assignmentString),
};

const REQUEST = {
    apm_query: (v: object) => assignmentString("apm_query", v),
    apm_stats_query: (v: object) => assignmentString("apm_stats_query", v),
    change_type: (v: object) => assignmentString("change_type", v),
    compare_to: (v: object) => assignmentString("compare_to", v),
    increase_good: (v: object) => assignmentString("increase_good", v),
    log_query: (v: object) =>
        block("log_query", v, (k1, v1) => convertFromDefinition(LOG_QUERY, k1, v1)),
    order_by: (v: object) => assignmentString("order_by", v),
    order_dir: (v: object) => assignmentString("order_dir", v),
    process_query: (v: object) => assignmentString("process_query", v),
    q: (v: object) => assignmentString("q", v),
    rum_query: (v: object) => assignmentString("rum_query", v),
    security_query: (v: object) => assignmentString("security_query", v),
    show_present: (v: object) => assignmentString("show_present", v),
    style: (v: object) => block("style", v, assignmentString),
    display_type: (v: object) => assignmentString("display_type", v),
    metadata: (v: object[]) => blockList(v, "metadata", assignmentString),
    network_query: (v: object) => assignmentString("network_query", v),
    on_right_yaxis: (v: object) => assignmentString("on_right_yaxis", v),
    aggregator: (v: object) => assignmentString("aggregator", v),
    alias: (v: object) => assignmentString("alias", v),
    cell_display_mode: (v: object) => assignmentString("cell_display_mode", v),
    conditional_formats: (v: object[]) => blockList(v, "conditional_formats", assignmentString),
    limit: (v: object) => assignmentString("limit", v),
    order: (v: object) => assignmentString("order", v),
    fill: (v: object) => block("fill", v, assignmentString),
};

const LOG_QUERY = {
    index: (v: object) => assignmentString("index", v),
    compute: (v: object) => block("compute_query", v, assignmentString),
    group_by: (v: object[]) => blockList(v, "group_by", (k1, v1) => convertFromDefinition(GROUP_BY, k1, v1)),
    multi_compute: (v: object[]) => blockList(v, "multi_compute", assignmentString),
    search: (v: any) => assignmentString("search_query", v.query),
    search_query: (v: object) => assignmentString("search_query", v),
};

const GROUP_BY = {
    facet: (v: object) => assignmentString("facet", v),
    limit: (v: object) => assignmentString("limit", v),
    sort: (v: object) => block("sort_query", v, assignmentString),
    sort_query: (v: object) => block("sort_query", v, assignmentString),
};

function convertSort(v: object) {
    return typeof v === "string"
        ? assignmentString("sort", v)
        : block("sort", v, assignmentString);
}

function convertWidgets(value: object[]) {
    return blockList(value, "widget", (k1, v1) => convertFromDefinition(WIDGET, k1, v1));
}

function convertRequests(value: object| object[]) {
    if (Array.isArray(value)) {
        return blockList(value, "request", (k, v) => convertFromDefinition(REQUEST, k, v));
    }
    return block("request", value, (k, v) => convertFromDefinition(REQUEST, k, v));
}

function widgetDefinition(contents: any) {
    let definitionType = contents.type === "slo" ? "service_level_objective" : contents.type;
    return block(`${definitionType}_definition`, contents, (k, v) =>
        convertFromDefinition(WIDGET_DEFINITION, k, v)
    );
}

export function generateDashboardTerraformCode(resourceName: string, dashboardData: object) {
    let result = "";
    Object.entries(dashboardData).forEach(([k, v]) => {
        result += convertFromDefinition(DASHBOARD, k, v);
    });
    return `resource "datadog_dashboard" "${resourceName}" {${result}\n}`;
}
