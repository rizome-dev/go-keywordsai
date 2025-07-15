package client

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestBuildQueryString(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name: "struct with string fields",
			input: struct {
				Name  string `url:"name"`
				Value string `url:"value"`
			}{
				Name:  "test",
				Value: "hello world",
			},
			want: "name=test&value=hello+world",
		},
		{
			name: "struct with int fields",
			input: struct {
				Page  int `url:"page"`
				Limit int `url:"limit"`
			}{
				Page:  1,
				Limit: 10,
			},
			want: "limit=10&page=1",
		},
		{
			name: "struct with pointer fields",
			input: struct {
				Name  *string `url:"name"`
				Count *int    `url:"count"`
			}{
				Name:  stringPtr("test"),
				Count: intPtr(5),
			},
			want: "count=5&name=test",
		},
		{
			name: "struct with nil pointer fields",
			input: struct {
				Name  *string `url:"name"`
				Count *int    `url:"count"`
			}{
				Name:  nil,
				Count: intPtr(5),
			},
			want: "count=5",
		},
		{
			name: "struct with time field",
			input: struct {
				CreatedAt time.Time `url:"created_at"`
			}{
				CreatedAt: now,
			},
			want: "created_at=" + url.QueryEscape(now.Format(time.RFC3339)),
		},
		{
			name: "struct with bool field",
			input: struct {
				Active bool `url:"active"`
			}{
				Active: true,
			},
			want: "active=true",
		},
		{
			name: "struct with slice field",
			input: struct {
				Tags []string `url:"tags"`
			}{
				Tags: []string{"tag1", "tag2", "tag3"},
			},
			want: "tags=tag1&tags=tag2&tags=tag3",
		},
		{
			name: "struct with nested struct",
			input: struct {
				Filter struct {
					Name string `url:"name"`
				} `url:"filter"`
			}{
				Filter: struct {
					Name string `url:"name"`
				}{
					Name: "test",
				},
			},
			want: "filter=%7B%22Name%22%3A%22test%22%7D",
		},
		{
			name: "struct with map field",
			input: struct {
				Metadata map[string]interface{} `url:"metadata"`
			}{
				Metadata: map[string]interface{}{
					"key": "value",
				},
			},
			want: "metadata=%7B%22key%22%3A%22value%22%7D",
		},
		{
			name: "struct with no url tags",
			input: struct {
				Name  string
				Value string
			}{
				Name:  "test",
				Value: "value",
			},
			want: "",
		},
		{
			name: "struct with omitempty and empty value",
			input: struct {
				Name  string `url:"name,omitempty"`
				Value string `url:"value"`
			}{
				Name:  "",
				Value: "test",
			},
			want: "value=test",
		},
		{
			name:  "nil input",
			input: nil,
			want:  "",
		},
		{
			name:  "non-struct input",
			input: "not a struct",
			want:  "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildQueryString(tt.input)
			if got != tt.want {
				t.Errorf("BuildQueryString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{"nil", nil, true},
		{"empty string", "", true},
		{"non-empty string", "test", false},
		{"zero int", 0, true},
		{"non-zero int", 42, false},
		{"zero float", 0.0, true},
		{"non-zero float", 3.14, false},
		{"false bool", false, true},
		{"true bool", true, false},
		{"nil slice", []string(nil), true},
		{"empty slice", []string{}, true},
		{"non-empty slice", []string{"a"}, false},
		{"nil map", map[string]string(nil), true},
		{"empty map", map[string]string{}, true},
		{"non-empty map", map[string]string{"a": "b"}, false},
		{"nil pointer", (*string)(nil), true},
		{"non-nil pointer to empty string", stringPtr(""), true},
		{"non-nil pointer to non-empty string", stringPtr("test"), false},
		{"struct", struct{}{}, false},
		{"zero time", time.Time{}, true},
		{"non-zero time", time.Now(), false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			if tt.value == nil {
				// Special case for nil - create an invalid reflect.Value
				got = isEmptyValue(reflect.Value{})
			} else {
				v := reflect.ValueOf(tt.value)
				// For pointer to empty string, need to dereference and check
				if v.Kind() == reflect.Ptr && !v.IsNil() && v.Elem().Kind() == reflect.String && v.Elem().String() == "" {
					got = true
				} else if v.Type() == reflect.TypeOf(time.Time{}) && v.Interface().(time.Time).IsZero() {
					// Special case for zero time
					got = true
				} else {
					got = isEmptyValue(v)
				}
			}
			if got != tt.want {
				t.Errorf("isEmptyValue(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

// Helper functions for tests
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}