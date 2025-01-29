package util

import (
	"regexp"
	"strings"
)

func IsEmptyString(str string) bool {
	return str == ""
}

func IsEmptyStringWithTrimSpace(str *string) bool {
	if *str = strings.TrimSpace(*str); IsEmptyString(*str) {
		return true
	}
	return false
}

func ValidateRegex(value, exp string) (bool, error) {
	rgx, err := regexp.Compile(exp)
	if err != nil {
		return false, err
	}
	return rgx.MatchString(value), nil
}
