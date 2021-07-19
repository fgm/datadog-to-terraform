package convert

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

const indentFormat = "  "

type DataDogDocumentType int

const (
	InvalidType DataDogDocumentType = iota
	DashboardType
	MonitorType
)

type JSONData = map[string]interface{}

func isSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

type Description struct {
	IsBlockList bool
	IsHidden    bool
	IsRequired  bool
	Converter
}

func indent(level int) string {
	if level < 0 {
		log.Printf("Indent called with negative level %d. Adjusting to 0.", level)
		level = 0
	}
	return strings.Repeat(indentFormat, level)
}

type TFBlock struct {
	JSONTag    string
	RootTag    string
	props      map[string]Renderable
	blockLists map[string]Renderable
}

func (b *TFBlock) loadBlockList(name string, v interface{}, d *Description) error {
	r, err := d.Converter(name, v)
	if err != nil {
		return fmt.Errorf("converting key %s in %s context %s: %w",
			name, b.RootTag, b.JSONTag, err)
	}
	b.blockLists[name] = r
	return nil
}

func (b *TFBlock) loadProp(k string, v interface{}, d *Description) error {
	r, err := d.Converter("", v)
	if err != nil {
		return fmt.Errorf("converting key %s in %s context %s: %w",
			k, b.RootTag, b.JSONTag, err)
	}
	b.props[k] = r
	return nil
}

func (b *TFBlock) Load(jd JSONData, schemaPath [2]string) error {
	schema := schemas(schemaPath[0], schemaPath[1])
	for k, v := range jd {
		itemSchema, ok := schema[k]
		if !ok {
			return fmt.Errorf("unexpected key %s in %s context %s",
				k, b.RootTag, b.JSONTag)
		}
		if !itemSchema.IsBlockList {
			b.loadProp(k, v, itemSchema)
			continue
		}
		b.loadBlockList(k, v, itemSchema)
	}
	return nil
}

type TFBlockList []TFBlock

func schemas(l1, l2 string) map[string]*Description {
	complete := map[string]map[string]map[string]*Description{
		// Keys are the names of the keys in the JSON version.
		// They usually do NOT match the TF block names, e.g "widgets" in JSON
		// translates to a list of "widget" blocks with equivalent content; while
		// "queries" in JSON translates to a list of "query" blocks each containing a
		// single "metric_query" block with equivalent content.
		"dashboard": {
			"dashboard": {
				// IsRequired
				"layout_type": {Converter: convertString, IsRequired: true},
				"title":       {Converter: convertString, IsRequired: true},
				"widgets":     {Converter: convertWidgetList, IsRequired: true, IsBlockList: true},

				// Normal
				"dashboard_lists":  {Converter: convertIntList},
				"description":      {Converter: convertString},
				"is_read_only":     {Converter: convertBool},
				"notify_list":      {Converter: convertStringList},
				"reflow_type":      {Converter: convertString}, // "auto"|"fixed". Only if layout_type is "ordered"
				"restricted_roles": {Converter: convertStringList},
				"url":              {Converter: convertString},

				// IsHidden
				"dashboard_lists_removed": {Converter: convertHidden, IsHidden: true},
				"id":                      {Converter: convertHidden, IsHidden: true},

				// TODO
				"template_variables":        {Converter: makeBlockVal("template_variable", nil), IsBlockList: true},
				"template_variable_presets": {Converter: makeBlockVal("template_variable_preset", nil), IsBlockList: true},
			},
			"widget": {
				"definition":    {Converter: convertWidgetDefinition, IsBlockList: true},
				"id":            {Converter: convertHidden, IsHidden: true},
				"widget_layout": {Converter: convertHidden, IsBlockList: true, IsHidden: true},
			},
			"definition": {

				"alert_id":            {Converter: convertString},
				"autoscale":           {Converter: z},
				"background_color":    {Converter: convertString},
				// "banner_img" in group_definition
				"check":               {Converter: convertString},
				"color":               {Converter: convertString},
				"color_by_groups":     {Converter: z},
				"color_preference":    {Converter: z},
				"columns":             {Converter: z},
				"content":             {Converter: z},
				"count":               {Converter: z},
				"custom_links":        {Converter: makeBlockVal("custom_link", nil)},
				"custom_unit":         {Converter: z},
				"display_format":      {Converter: z},
				"env":                 {Converter: z},
				"event":               {Converter: makeBlockVal("event", nil)},
				"event_size":          {Converter: convertString},
				"filters":             {Converter: z},
				"font_size":           {Converter: convertString},
				"global_time_target":  {Converter: z},
				"group":               {Converter: convertString},
				"group_by":            {Converter: convertStringList},
				"grouping":            {Converter: convertString},
				"has_padding":         {Converter: z},
				"has_search_bar":      {Converter: z},
				"hide_zero_counts":    {Converter: z},
				"indexes":             {Converter: z},
				"layout_type":         {Converter: convertString},
				"legend_columns":      {Converter: z},
				"legend_layout":       {Converter: z},
				"legend_size":         {Converter: convertString},
				"live_span":           {Converter: convertString},
				"logset":              {Converter: z},
				"margin":              {Converter: z},
				"markers":             {Converter: z},
				"message_display":     {Converter: z},
				"no_group_hosts":      {Converter: z},
				"no_metric_hosts":     {Converter: z},
				"node_type":           {Converter: z},
				"precision":           {Converter: convertInt},
				"query":               {Converter: convertString},
				"requests":            {Converter: makeBlockVal("request", nil)},
				"right_yaxis":         {Converter: z},
				"scope":               {Converter: z},
				"service":             {Converter: z},
				"show_breakdown":      {Converter: z},
				"show_date_column":    {Converter: z},
				"show_distribution":   {Converter: z},
				"show_error_budget":   {Converter: z},
				"show_errors":         {Converter: z},
				"show_hits":           {Converter: z},
				"show_last_triggered": {Converter: z},
				"show_latency":        {Converter: z},
				"show_legend":         {Converter: convertBool},
				"show_message_column": {Converter: z},
				"show_resource_list":  {Converter: z},
				"show_tick":           {Converter: z},
				// "show_title" in grou_definition
				"size_format":         {Converter: z},
				"sizing":              {Converter: z},
				"slo_id":              {Converter: z},
				"sort":                {Converter: z},
				"span_name":           {Converter: z},
				"start":               {Converter: z},
				"style":               {Converter: makeBlockVal("style", nil)},
				"summary_type":        {Converter: z},
				"tags":                {Converter: convertStringList},
				"tags_execution":      {Converter: convertString},
				"text":                {Converter: convertString},
				"text_align":          {Converter: convertString},
				"tick_edge":           {Converter: z},
				"tick_pos":            {Converter: z},
				"time":                {Converter: z},
				"time_windows":        {Converter: z},
				"title":               {Converter: convertString},
				"title_align":         {Converter: convertString},
				"title_size":          {Converter: convertString},
				"type":                {Converter: z},
				"unit":                {Converter: convertString},
				"url":                 {Converter: z},
				"vertical_align":      {Converter: z},
				"view_mode":           {Converter: z},
				// "view" ? in geomap_definition
				"view_type":           {Converter: z},
				"viz_type":            {Converter: convertString},
				"widget_layout":       {Converter: z},
				"widgets":             {Converter: convertWidgetList},
				"xaxis":               {Converter: z},
				"yaxis":               {Converter: z},
			},
		},

		"monitor": {
			"monitor": {},
		},
	}
	return complete[l1][l2]
}
