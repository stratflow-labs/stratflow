package testkit

//go:fix inline
func StringPtr(value string) *string {
	return new(value)
}

//go:fix inline
func Int32Ptr(value int32) *int32 {
	return new(value)
}
