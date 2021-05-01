package main

import (
	"database/sql"
	"testing"
)

/*** Integration Tests ***/
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
			0,
			ksatNameErr{name: ""},
		},
		{
			"OneLetterName",
			"a",
			0,
			noKsatIdByNameErr{name: "a", err: sql.ErrNoRows},
		},
		{
			"MultipleLetterName",
			"abc",
			0,
			noKsatIdByNameErr{name: "abc", err: sql.ErrNoRows},
		},
		{
			"MultipleLetterNameWhichExists (first)",
			"first", // in db as a ksat (first one created)
			1,
			nil,
		},
		{
			"MultipleLetterNameWhichExists (second)",
			"second", // in db as a ksat (second one created)
			2,
			nil,
		},
		{
			"SixLetterName",
			"sixsix",
			0,
			noKsatIdByNameErr{name: "sixsix", err: sql.ErrNoRows},
		},
		{
			"SevenLetterName",
			"sevsevs",
			0,
			ksatNameErr{name: "sevsevs"},
		},
		{
			"NameWithNumberAndLetters",
			"a12b",
			0,
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
