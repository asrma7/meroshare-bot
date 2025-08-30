package utils

import (
	"strconv"
	"strings"
)

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(strings.ToLower(s[:1])) + s[1:]
}

func StringToInt(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}
