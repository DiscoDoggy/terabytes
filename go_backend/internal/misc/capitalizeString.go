package misc

import "strings"

func CapitalizeString(str string) string {
	if len(str) == 0 {
		return ""
	}

	return strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
}