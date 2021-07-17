package convert

import (
	"reflect"
)

const indent = "  "

type DataDogDocumentType int

const (
	InvalidType DataDogDocumentType = iota
	DashboardType
	MonitorType
)

type JSONData map[string]interface{}

type Converter func(interface{}) string

type Convertible interface {
	IndentedString(indent int) string
}

type Block []JSONData

type Prop JSONData

type Document struct {
	tag    string
	blocks map[string]Block
	props  map[string]Prop
}

func (d *Document) String() string {
	return ""
}

func isSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}
