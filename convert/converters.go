package convert

import (
	"errors"
	"fmt"
)

// Converter implementations transforms a value of an expected concrete type to a Renderable.
type Converter func(string, interface{}) (Renderable, error)

func convertWidgetDefinition(name string, v interface{}) (Renderable, error) {
	const root = "dashboard"
	schema := schemas(root, "widget")
	b := TFBlock{
		JSONTag:    "definition",
		RootTag:    root,
		props:      nil,
		blockLists: nil,
	}
	return nil, nil
}

func convertWidgetList(name string, v interface{}) (Renderable, error) {
	switch cast := v.(type) {
	case nil:
		return TFBlockList{}, nil
	case []interface{}:
		if len(cast) == 0 {
			return TFBlockList{}, nil
		}
		bl := make(TFBlockList, len(cast))
		for i, item := range cast {
			child := TFBlock{
				JSONTag:    "widget",
				RootTag:    "dashboard",
				blockLists: make(map[string]Renderable),
				props:      make(map[string]Renderable),
			}
			mi, ok := item.(JSONData)
			if !ok {
				return nil, fmt.Errorf("convertWidgetList could not convert %#v to JSONData", v)
			}
			if err := child.Load(mi, [2]string{child.RootTag, child.JSONTag}); err != nil {
				return nil, fmt.Errorf("loading %s bl: %w", name, err)
			}
			bl[i] = child
		}
		return TFBlockList{}, errors.New("convertIntList cannot convert []interface{}")
	default:
		return TFBlockList{}, nil
	}

	return nil, nil
}

func convertBool(_ string, v interface{}) (Renderable, error) {
	bo, ok := v.(bool)
	if !ok {
		return boolValue(false), fmt.Errorf("convertBool cannot convert %T", v)
	}
	return boolValue(bo), nil
}

func convertHidden(_ string, _ interface{}) (Renderable, error) {
	return hiddenValue(""), nil
}

func convertInt(_ string, v interface{})  (Renderable, error) {
	iv, ok := v.(int)
	if !ok {
		return intValue(0), fmt.Errorf("convertInt cannot convert %T", v)
	}
	return intValue(iv), nil
}

func convertIntList(tag string, v interface{}) (Renderable, error) {
	switch cast := v.(type) {
	case nil:
		return intListValue{}, nil
	case []interface{}:
		if len(cast) == 0 {
			return stringValue("[]"), nil
		}
		return intListValue{}, errors.New("convertIntList cannot convert []interface{}")
	case []int:
		return intListValue(cast), nil
	default:
		return intListValue{}, fmt.Errorf("convertIntList cannot convert type: %T", cast)
	}
}

func convertStringList(tag string, v interface{}) (Renderable, error) {
	switch cast := v.(type) {
	case nil:
		return stringListValue{}, nil
	case []interface{}:
		if len(cast) == 0 {
			return stringListValue{}, nil
		}
		return stringListValue{}, errors.New("convertStringList cannot convert []interface{}")
	case []string:
		return stringListValue(cast), nil
	default:
		return stringListValue{}, fmt.Errorf("convertStringList cannot %T", cast)
	}
}

func convertString(_ string, v interface{}) (Renderable, error) {
	sv, ok := v.(string)
	if !ok {
		return stringValue(""), fmt.Errorf("convertString cannot convert %T", v)
	}
	return stringValue(sv), nil
}
