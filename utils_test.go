package main

import (
	"testing"
)

func TestProfanityFilter(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"This is a kerfuffle", "This is a ****"},
		{"Sharbert is a word", "**** is a word"},
		{"Fornax is banned", "**** is banned"},
		{"No banned words here", "No banned words here"},
		{"Mixed case Kerfuffle", "Mixed case ****"},
	}

	for _, test := range tests {
		result := profanityFilter(test.input)
		if result != test.expected {
			t.Errorf("Expected '%s', but got '%s'", test.expected, result)
		}
	}
}