package utils

func FormatStringPtr(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
