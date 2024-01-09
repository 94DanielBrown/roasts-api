package utils

import (
	"testing"
)

func TestToPascalCase(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"SingleWord", "coffee", "Coffee"},
		{"MultipleWords", "cranberry greens powder shake", "CranberryGreensPowderShake"},
		{"MixedCase", "Peppermint Tea", "PeppermintTea"},
		{"Spaces", "  water  ", "Water"},
		{"EmptyString", "", ""},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToPascalCase(tc.input)
			if result != tc.expected {
				t.Errorf("ToPascalCase(%v) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}
