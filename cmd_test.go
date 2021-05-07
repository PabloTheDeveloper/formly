package main

import (
	"database/sql"
	"testing"
)

/*** Integration Tests ***/
func TestNewCommand(t *testing.T) {
	type expected struct {
		command *command
		error
	}
	cases := []struct {
		desc string
		args []string
		expected
	}{
		{
			"valid command name (exist in db)", []string{"ksat", "first"}, expected{
				&command{ksat: ksat{id: 1}},
				nil,
			},
		},
		{
			"valid command name (dne in db)", []string{"ksat", "zzero"}, expected{nil, sql.ErrNoRows},
		},
		{
			"invalid command name (dne in db)", []string{"ksat", "z z"}, expected{nil, wordErr{"z z"}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			given := expected{}
			if cmd, err := newCommand(tc.args); cmd != nil || err != nil {
				given.command = cmd
				given.error = err
			}
			if given.error != tc.expected.error {
				t.Fatalf("given error: %v\nexpected error: %v", given.error, tc.expected.error)
			}
			if given.command == nil && tc.expected.command == nil {
				return
			}
			if given.command == nil || tc.expected.command == nil {
				t.Fatalf(
					"given command: %v\nexpected command: %v",
					given.command,
					tc.expected.command,
				)
			}
			if given.command.ksat.id != tc.expected.command.ksat.id {
				t.Fatalf(
					"given command.ksat.id: %v\nexpected command.ksat.id: %v",
					given.command.ksat.id,
					tc.expected.command.ksat.id,
				)
			}
		})
	}
}
