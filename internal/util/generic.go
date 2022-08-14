package util

// Ptr returns a pointer to a value
func Ptr[T any](v T) *T {
	return &v
}
