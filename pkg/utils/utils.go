package utils

import (
	"strings"
	"unicode"
)

func ToPascalCase(s string) string {
	var pascalCase strings.Builder
	nextToUpper := true

	s = strings.TrimSpace(s)

	for _, r := range s {
		if nextToUpper {
			pascalCase.WriteRune(unicode.ToUpper(r))
			nextToUpper = false
		} else if r == ' ' {
			nextToUpper = true
		} else {
			pascalCase.WriteRune(r)
		}
	}

	return pascalCase.String()
}
