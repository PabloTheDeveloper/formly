package main

import "testing"

/*** Unit Tests ***/
func TestIsWordValid(t *testing.T) {
	cases := []struct {
		desc     string
		word     string
		expected error
	}{
		{
			"ZeroLetterWord",
			"",
			wordErr{word: ""},
		},
		{
			"OneLetterWord",
			"a",
			nil,
		},
		{
			"MultipleLetterWord",
			"abc",
			nil,
		},
		{
			"WordWithNumberAndLetters",
			"a12b",
			wordErr{word: "a12b"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			if ret := isWordValid(tc.word); ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
		})
	}
}

func TestIsStringLengthCorrect(t *testing.T) {
	// Test not Needed.
}
