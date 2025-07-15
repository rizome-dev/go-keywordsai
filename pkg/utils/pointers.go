package utils

// String returns a pointer to the string value
func String(s string) *string {
	return &s
}

// Int returns a pointer to the int value
func Int(i int) *int {
	return &i
}

// Int64 returns a pointer to the int64 value
func Int64(i int64) *int64 {
	return &i
}

// Float64 returns a pointer to the float64 value
func Float64(f float64) *float64 {
	return &f
}

// Bool returns a pointer to the bool value
func Bool(b bool) *bool {
	return &b
}