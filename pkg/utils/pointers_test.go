package utils

import (
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"simple string", "hello"},
		{"string with spaces", "hello world"},
		{"unicode string", "こんにちは"},
		{"special characters", "!@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := String(tt.input)
			if result == nil {
				t.Fatal("expected non-nil pointer")
			}
			if *result != tt.input {
				t.Errorf("String() = %v, want %v", *result, tt.input)
			}
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"max int", int(^uint(0) >> 1)},
		{"min int", -int(^uint(0)>>1) - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int(tt.input)
			if result == nil {
				t.Fatal("expected non-nil pointer")
			}
			if *result != tt.input {
				t.Errorf("Int() = %v, want %v", *result, tt.input)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		name  string
		input int64
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"large positive", 9223372036854775807},
		{"large negative", -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int64(tt.input)
			if result == nil {
				t.Fatal("expected non-nil pointer")
			}
			if *result != tt.input {
				t.Errorf("Int64() = %v, want %v", *result, tt.input)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		name  string
		input float64
	}{
		{"zero", 0.0},
		{"positive", 42.5},
		{"negative", -42.5},
		{"very small", 0.0000001},
		{"very large", 1e10},
		{"pi", 3.14159265359},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Float64(tt.input)
			if result == nil {
				t.Fatal("expected non-nil pointer")
			}
			if *result != tt.input {
				t.Errorf("Float64() = %v, want %v", *result, tt.input)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name  string
		input bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Bool(tt.input)
			if result == nil {
				t.Fatal("expected non-nil pointer")
			}
			if *result != tt.input {
				t.Errorf("Bool() = %v, want %v", *result, tt.input)
			}
		})
	}
}

// Benchmark tests
func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = String("test")
	}
}

func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Int(42)
	}
}

func BenchmarkInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Int64(42)
	}
}

func BenchmarkFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Float64(3.14)
	}
}

func BenchmarkBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Bool(true)
	}
}