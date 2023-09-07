package ethereum_client

import (
	"strconv"
	"testing"
)

func TestBaseDecimalToHexString(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0"},
		{10, "a"},
		{255, "ff"},
		{-1, "-1"},
	}

	for _, test := range tests {
		t.Run(strconv.FormatInt(test.input, 10), func(t *testing.T) {
			result := baseDecimalToHexString(test.input)
			if result != test.expected {
				t.Errorf("Expected %s, but got %s", test.expected, result)
			}
		})
	}
}

func TestHexStringToDecimal(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0", 0},
		{"a", 10},
		{"ff", 255},
		{"-1", -1},
		{"invalid", 0},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := hexStringToDecimal(test.input)

			if err != nil {
				// Expect an error for invalid input.
				if test.expected != 0 {
					t.Errorf("Expected error, but got value %d", result)
				}
			} else if result != test.expected {
				t.Errorf("Expected %d, but got %d", test.expected, result)
			}
		})
	}
}
