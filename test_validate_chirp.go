package main

import "testing"

func testValidateChirp(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "This is a kerfuffle",
			expected: "This is a ****",
		},
	}

	for _, cas := range cases {
		actual, err := validate_chirp(cas.input)
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
		if actual != cas.expected {
			t.Errorf("String should be %s not %s", cas.expected, actual)
		}
	}
}
