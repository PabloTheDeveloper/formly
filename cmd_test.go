package main

import (
	"database/sql"
	"flag"
	"fmt"
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
				&command{Ksat: Ksat{id: 1}},
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
			if given.command.Ksat.id != tc.expected.command.Ksat.id {
				t.Fatalf(
					"given command.Ksat.id: %v\nexpected command.Ksat.id: %v",
					given.command.Ksat.id,
					tc.expected.command.Ksat.id,
				)
			}
		})
	}
}
func TestExecuteCommand(t *testing.T) {
	type expected struct {
		command command
		error
	}
	cases := []struct {
		desc string
		command
		expected
	}{
		{
			"valid prompts and flags set (and exist in db)",
			command{
				Ksat:     Ksat{id: 3},
				FlagSet:  flag.NewFlagSet("hasP", flag.ExitOnError),
				flagArgs: []string{"firstflag='data'", "secondflag='data'"},
			},
			expected{
				command{},
				nil,
			},
		},
		{
			"valid prompts but missing flag (and exist in db) (should work)",
			command{
				Ksat:     Ksat{id: 3},
				FlagSet:  flag.NewFlagSet("hasP", flag.ExitOnError),
				flagArgs: []string{"firstflag='data'"},
			},
			expected{
				command{},
				nil,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			given := expected{}
			given.error = tc.command.executeCommand()
			if given.error != tc.expected.error {
				t.Fatalf("given error: %v\nexpected error: %v", given.error, tc.expected.error)
			}
		})
		fmt.Println(tc.command.entries)
	}
}
