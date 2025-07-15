package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// BuildQueryString builds a URL query string from a struct using url tags
func BuildQueryString(params interface{}) string {
	if params == nil {
		return ""
	}

	values := url.Values{}
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Get the url tag
		urlTag := field.Tag.Get("url")
		if urlTag == "" || urlTag == "-" {
			continue
		}

		// Parse the url tag
		tagParts := strings.Split(urlTag, ",")
		fieldName := tagParts[0]

		// Skip if omitempty and value is empty
		if len(tagParts) > 1 && tagParts[1] == "omitempty" && isEmptyValue(value) {
			continue
		}

		// Handle different types
		if value.Kind() == reflect.Ptr && !value.IsNil() {
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.String:
			values.Set(fieldName, value.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values.Set(fieldName, fmt.Sprintf("%d", value.Int()))
		case reflect.Bool:
			values.Set(fieldName, fmt.Sprintf("%t", value.Bool()))
		case reflect.Float32, reflect.Float64:
			values.Set(fieldName, fmt.Sprintf("%f", value.Float()))
		case reflect.Slice:
			for j := 0; j < value.Len(); j++ {
				values.Add(fieldName, fmt.Sprintf("%v", value.Index(j)))
			}
		default:
			// Handle time.Time
			if value.Type() == reflect.TypeOf(time.Time{}) {
				t := value.Interface().(time.Time)
				values.Set(fieldName, t.Format(time.RFC3339))
			} else if value.Kind() == reflect.Struct || value.Kind() == reflect.Map {
				// For complex types, encode as JSON
				if jsonBytes, err := json.Marshal(value.Interface()); err == nil {
					values.Set(fieldName, string(jsonBytes))
				}
			}
		}
	}

	return values.Encode()
}

func isEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}