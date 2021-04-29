package main

import (
	"errors"
	"testing"
)

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
			err := newKsat(tc.name, tc.usage, tc.prompts)

			if nameBoundsErr != err {
				t.Fatalf("errors don't match, %v, %v", err, tc.expected)
			}
		})
	}
}
