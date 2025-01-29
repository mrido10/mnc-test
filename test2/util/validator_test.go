package util

import "testing"

func TestIsEmptyString(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{" ", false},
		{"test", false},
	}

	for _, test := range tests {
		if got := IsEmptyString(test.input); got != test.expected {
			t.Errorf("IsEmptyString(%q) = %v; want %v", test.input, got, test.expected)
		}
	}
}

func TestIsEmptyStringWithTrimSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{" ", true},
		{"  test  ", false},
	}

	for _, test := range tests {
		input := test.input
		if got := IsEmptyStringWithTrimSpace(&input); got != test.expected {
			t.Errorf("IsEmptyStringWithTrimSpace(%q) = %v; want %v", test.input, got, test.expected)
		}
	}
}

func TestValidateRegex(t *testing.T) {
	tests := []struct {
		value    string
		exp      string
		expected bool
		err      bool
	}{
		{"085289764830", "^(\\+62|62|0)8[1-9][0-9]{6,10}$", true, false},
		{"+6285289764830", "^(\\+62|62|0)8[1-9][0-9]{6,10}$", true, false},
		{"6285289764830", "^(\\+62|62|0)8[1-9][0-9]{6,10}$", true, false},
	}

	for _, test := range tests {
		got, err := ValidateRegex(test.value, test.exp)
		if (err != nil) != test.err {
			t.Errorf("ValidateRegex(%q, %q) unexpected error: %v", test.value, test.exp, err)
		}
		if got != test.expected {
			t.Errorf("ValidateRegex(%q, %q) = %v; want %v", test.value, test.exp, got, test.expected)
		}
	}
}
