package utils

func StrDefaultIfEmpty(str, defaultStr string) string {
	if str == "" {
		return defaultStr
	}
	return str
}
