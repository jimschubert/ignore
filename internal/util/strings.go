package util

// StringValue returns de-referenced string or empty string.
func StringValue(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
