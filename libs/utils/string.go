package utils

import "strings"

func NormalizeString(str string) string {
	str = strings.ToLower(str)
	return strings.TrimSpace(str)
}
