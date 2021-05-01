package main

import "testing"

func TestGetKsatIdByName(t *testing.T) {
	cases := []struct {
		desc        string
		name        string
		expectedId  int64
		expectedErr error
	}{
		{
			"ZeroLetterName",
			"",
			-1,
			ksatNameErr{name: ""},
		},
		{
			"OneLetterName",
			"a",
			-1,
			nil,
		},
		{
			"MultipleLetterName",
			"abc",
			-1,
			nil,
		},
		{
			"MultipleLetterNameWhichExists",
			"first", // in db as a ksat (first one created)
			1,
			nil,
		},
		{
			"MultipleLetterNameWhichExists",
			"second", // in db as a ksat (second one created)
			2,
			nil,
		},
		{
			"SixLetterName",
			"sixsix",
			-1,
			nil,
		},
		{
			"SevenLetterName",
			"sevsevs",
			-1,
			ksatNameErr{name: "sevsevs"},
		},
		{
			"NameWithNumberAndLetters",
			"a12b",
			-1,
			wordErr{word: "a12b"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			id, err := getKsatIdByName(tc.name)
			if id != tc.expectedId {
				t.Fatalf("ids don't match: %v, %v", id, tc.expectedId)
			}
			if err != tc.expectedErr {
				t.Fatalf("errors don't match: %v, %v", err, tc.expectedErr)
			}
		})
	}
}

/*
func TestNewKsat(t *testing.T) {
	nameBoundsErr := errors.New("'name' must be between 1-6 characters long")
	validUsageStr := "some usage"
	validPrompt := "some prompt"
	newKsatCases := []struct {
		desc                 string
		name, usage, prompts string
		expected             error
	}{
		{
			"testing invalid ksat name (too low)",
			"",
			validUsageStr,
			validPrompt,
			nameBoundsErr,
		},
		{
			"testing valid ksat name (lowest inclusion point)",
			"o",
			validUsageStr,
			validPrompt,
			nil,
		},
		{
			"testing valid ksat name (highest inclusion point)",
			"onesix",
			validUsageStr,
			validPrompt,
			nil,
		},
		{
			"testing invalid ksat name (too high)",
			"onesix7",
			validUsageStr,
			validPrompt,
			nameBoundsErr,
		},
	}

	for _, tc := range newKsatCases {
		t.Run(tc.desc, func(t *testing.T) {
			// err := newKsat(tc.name, tc.usage, tc.prompts)

			// if nameBoundsErr != err {
			// t.Fatalf("errors don't match, %v, %v", err, tc.expected)
			// }
		})
	}
}

*/
