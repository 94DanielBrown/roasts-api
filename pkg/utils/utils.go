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

// CalculateAverageRating takes a slice of float64 values (ratings) and returns the average
func CalculateAverageRating(ratings []float64) float64 {
	if len(ratings) == 0 {
		return 0
	}

	var sum float64
	for _, rating := range ratings {
		sum += rating
	}

	return sum / float64(len(ratings))
}
