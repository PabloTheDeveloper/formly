package main

import (
	"database/sql"
	"math"
	"reflect"
	"testing"
	"time"
)

/*** Unit Tests ***/
func TestSetName(t *testing.T) {
	type output struct {
		Ksat
		error
	}
	cases := []struct {
		desc     string
		arg1     string
		expected output
	}{
		{"ZeroLetterWord", "", output{Ksat{}, strLengthErr{lower: 1, upper: 6, str: ""}}},
		{"OneLetterWord", "second", output{Ksat{name: "second"}, nil}},
		{"SixLetterWord", "sixsix", output{Ksat{name: "sixsix"}, nil}},
		{"SevenLetterWord", "sevsevs", output{Ksat{},
			strLengthErr{lower: 1, upper: 6, str: "sevsevs"}}},
		{"WordWithNumberAndLetters", "a12b", output{Ksat{}, wordErr{word: "a12b"}}},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			ksat := Ksat{}
			err := ksat.SetName(tc.arg1)
			returned := output{ksat, err}
			if returned.error != tc.expected.error {
				t.Fatalf("errors don't match\nreturned err: %v\nexpected err: %v", returned.error, tc.expected)
			}
			if !reflect.DeepEqual(returned.Ksat, tc.expected.Ksat) {
				t.Fatalf("ksats don't match\nreturned ksat: %v\nexpected ksat: %v", returned.Ksat, tc.expected.Ksat)
			}
		})
	}
}
func TestSetUsage(t *testing.T) {
	type output struct {
		Ksat
		error
	}
	cases := []struct {
		desc     string
		arg1     string
		expected output
	}{
		{"FourLetterUsage", "abcd",
			output{Ksat{}, strLengthErr{lower: 5, upper: 40, str: "abcd"}}},
		{"FiveLetterUsage", "abcde", output{Ksat{}, nil}},
		{"FourtyLetterUsage", "0123456789012345678901234567890123456789", output{Ksat{}, nil}},
		{"Fourty1LetterUsage", "0123456789012345678901234567890123456789" + "1",
			output{Ksat{}, strLengthErr{lower: 5, upper: 40, str: "0123456789012345678901234567890123456789" + "1"}}},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			ksat := Ksat{}
			err := ksat.SetUsage(tc.arg1)
			returned := output{ksat, err}
			if returned.error != tc.expected.error {
				t.Fatalf("errors don't match\nreturned err: %v\nexpected err: %v", returned.error, tc.expected)
			}
			if !reflect.DeepEqual(returned.Ksat, tc.expected.Ksat) {
				t.Fatalf("ksats don't match\nreturned ksat: %v\nexpected ksat: %v", returned.Ksat, tc.expected.Ksat)
			}
		})
	}
}

/*** Integration Tests ***/
func TestGetKsatByName(t *testing.T) {
	type output struct {
		Ksat
		error
	}
	cases := []struct {
		desc     string
		arg1     string
		expected output
	}{
		{
			"valid name but for a Ksat that exists",
			"first",
			output{
				Ksat{id: 1, name: "first", usage: "some usage"},
				nil,
			},
		},
		{
			"another valid name but for a Ksat that exists",
			"second",
			output{
				Ksat{id: 2, name: "second", usage: "some more usage"},
				nil,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			ksat, err := GetKsatByName(tc.arg1)
			returned := output{ksat, err}
			if returned.error != tc.expected.error {
				t.Fatalf("errors don't match\nreturned err: %v\nexpected err: %v", returned.error, tc.expected)
			}
			if !reflect.DeepEqual(returned.Ksat, tc.expected.Ksat) {
				t.Fatalf("ksats don't match\nreturned ksat: %v\nexpected ksat: %v", returned.Ksat, tc.expected.Ksat)
			}
		})
	}
}
func TestGetKsatByID(t *testing.T) {
	type expected struct {
		Ksat
		error
	}
	taskCases := []struct {
		desc   string
		ksatID int64
		expected
	}{
		{
			"valid id for a Ksat that exists",
			1,
			expected{
				Ksat:  Ksat{id: 1, name: "first", usage: "some usage"},
				error: nil,
			},
		},
		{
			"valid id for a Ksat that does not exist",
			1000,
			expected{
				Ksat:  Ksat{},
				error: sql.ErrNoRows,
			},
		},
	}
	for _, tc := range taskCases {
		t.Run(tc.desc, func(t *testing.T) {
			ksat, err := GetKsatByID(tc.ksatID)
			returned := expected{ksat, err}
			if returned.error != tc.expected.error {
				t.Fatalf("errors don't match\nreturned err: %v\nexpected err: %v", returned.error, tc.expected)
			}
			if !reflect.DeepEqual(returned.Ksat, tc.expected.Ksat) {
				t.Fatalf("ksats don't match\nreturned ksat: %v\nexpected ksat: %v", returned.Ksat, tc.expected.Ksat)
			}
		})
	}
}
func TestGetByID(t *testing.T) {
	promptCases := []struct {
		desc       string
		prompt     prompt
		successful prompt
		expected   error
	}{
		{
			"valid id for a prompt that exists",
			prompt{
				id:       1,
				KsatID:   3,
				sequence: 1,
				flag:     "firstflag",
				usage:    "some usage",
			},
			prompt{
				id:       1,
				KsatID:   3,
				sequence: 1,
				flag:     "firstflag",
				usage:    "some usage",
			},
			nil,
		},
		{
			"valid id for a prompt that does not exist",
			prompt{id: 1000, sequence: 10, flag: "dne", usage: "second usage here"},
			prompt{id: 1000, sequence: 10, flag: "dne", usage: "second usage here"}, // needs to be same even if its suppose to fail
			sql.ErrNoRows,
		},
	}
	for _, tc := range promptCases {
		t.Run(tc.desc, func(t *testing.T) {
			ret := tc.prompt.getByID()
			if ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
			if tc.prompt.id != tc.successful.id {
				t.Fatalf("ids don't match: %v, %v", tc.prompt.id, tc.successful.id)
			}
			if tc.prompt.sequence != tc.successful.sequence {
				t.Fatalf("sequences don't match: %v, %v",
					tc.prompt.sequence, tc.successful.sequence)
			}
			if tc.prompt.flag != tc.successful.flag {
				t.Fatalf("flags don't match: %v, %v", tc.prompt.flag, tc.successful.flag)
			}
			if tc.prompt.usage != tc.successful.usage {
				t.Fatalf("usages don't match: %v, %v", tc.prompt.usage, tc.successful.usage)
			}
		})
	}
	sessionCases := []struct {
		desc       string
		session    session
		successful session
		expected   error
	}{
		{
			"valid id for a session that exists",
			session{
				id:       1,
				KsatID:   3,
				createAt: time.Date(2001, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			session{
				id:       1,
				KsatID:   3,
				createAt: time.Date(2001, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			nil,
		},
		{
			"valid id for a session that does not exist",
			session{
				id:       10101,
				KsatID:   3,
				createAt: time.Date(2001, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			session{
				id:       10101,
				KsatID:   3,
				createAt: time.Date(2001, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			sql.ErrNoRows,
		},
	}
	for _, tc := range sessionCases {
		t.Run(tc.desc, func(t *testing.T) {
			ret := tc.session.getByID()
			if ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
			if tc.session.id != tc.successful.id {
				t.Fatalf("ids don't match: %v, %v", tc.session.id, tc.successful.id)
			}
			if tc.session.KsatID != tc.successful.KsatID {
				t.Fatalf("KsatIDs don't match: %v, %v", tc.session.KsatID, tc.successful.KsatID)
			}
		})
	}
}
func TestGetPromptsByKsatID(t *testing.T) {
	type output struct {
		prompts []prompt
		error
	}
	cases := []struct {
		desc     string
		ksatID   int64
		expected output
	}{
		{
			"valid ID which contains 1 valid prompt",
			3,
			output{
				[]prompt{
					{id: 1, KsatID: 3, sequence: 1, flag: "firstflag", usage: "some usage"},
					{id: 2, KsatID: 3, sequence: 2, flag: "secondflag", usage: "some usage"},
				},
				nil,
			},
		},
		{
			"valid ID but it has no prompt",
			1,
			output{
				[]prompt{},
				nil,
			},
		},
		{
			"valid id for a Ksat that does not exist",
			1000,
			output{
				nil,
				sql.ErrNoRows,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			prompts, err := GetPromptsByKsatID(tc.ksatID)
			returned := output{prompts, err}
			if returned.error != tc.expected.error {
				t.Fatalf("errors don't match\nreturned err: %v\nexpected err: %v", returned.error, tc.expected)
			}

			if !reflect.DeepEqual(returned.prompts, tc.expected.prompts) {
				t.Fatalf("prompts don't match\nreturned ksat: %v\nexpected ksat: %v", returned.prompts, tc.expected.prompts)
			}
		})
	}
}
func TestGetSessionsByID(t *testing.T) {
	cases := []struct {
		desc     string
		ksatID   int64
		sessions []session
		err      error
	}{
		{
			"valid ID which contains 1 valid session",
			3,
			[]session{
				{
					id:       1,
					KsatID:   3,
					createAt: time.Date(2000, 11, 17, 20, 34, 58, 651387237, time.UTC),
				},
				{
					id:       2,
					KsatID:   3,
					createAt: time.Date(2001, 11, 17, 20, 34, 58, 651387237, time.UTC),
				},
			},
			nil,
		},
		{
			"valid ID but it has no session",
			1,
			[]session{},
			nil,
		},
		{
			"valid id for a Ksat that does not exist",
			1000,
			nil,
			sql.ErrNoRows,
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			sessions, err := GetSessionsByKsatID(tc.ksatID)
			if err != tc.err {
				t.Fatalf("errors don't match: %v, %v", err, tc.err)
			}
			for i, item := range sessions {
				// needed since I am creating more sessions than are shown
				if i >= len(tc.sessions) {
					break
				}
				if item.id != tc.sessions[i].id {
					t.Fatalf("ids don't match: %v, %v", item.id, tc.sessions[i].id)
				}
				if item.createAt != tc.sessions[i].createAt {
					t.Fatalf("createAts don't match: %v, %v", item.createAt, tc.sessions[i].createAt)
				}
			}
		})
	}
}
func TestDbInsert(t *testing.T) {
	KsatCases := []struct {
		desc             string
		task             Ksat
		successfulInsert bool
		expected         error
	}{
		{
			"invalid name for Ksat (letter + number)",
			Ksat{name: "a12b", usage: "some usage"},
			false,
			wordErr{word: "a12b"},
		},
		{
			"valid name but for a Ksat that exists",
			Ksat{name: "first", usage: "first usage here"},
			false,
			AlreadyExistsErr{identifier: "first", tableName: "Ksat"},
		},
		{
			"valid name but for a Ksat that does not exists",
			Ksat{name: "new", usage: "usage here"},
			true,
			nil,
		},
	}
	for _, tc := range KsatCases {
		t.Run(tc.desc, func(t *testing.T) {
			ret := tc.task.dbInsert()
			if ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
			if tc.successfulInsert && tc.task.id == 0 {
				t.Fatalf("error. Id for new Ksat is not assigned")
			}
		})
	}
	promptCases := []struct {
		desc             string
		prompt           prompt
		successfulInsert bool
		expected         error
	}{
		{
			"perfect prompt but the Ksat_id makes it invalid",
			prompt{KsatID: 0, sequence: 1, flag: "firstflag", usage: "some usage"},
			false,
			sql.ErrNoRows,
		},
		{
			"perfect prompt but the sequence makes it invalid",
			prompt{KsatID: 1, sequence: -1, flag: "firstflag", usage: "some usage"},
			false,
			numLengthErr{lower: 1, upper: math.MaxInt64, num: -1},
		},
		{
			"valid prompt (no sequence conflict since it is the first flag for 'first' Ksat)",
			prompt{KsatID: 1, sequence: 1, flag: "firstflag", usage: "some usage"},
			true,
			nil,
		},
		{
			"valid prompt (no sequence conflict)",
			prompt{KsatID: 1, sequence: 2, flag: "secondflag", usage: "some usage"},
			true,
			nil,
		},
	}
	for _, tc := range promptCases {
		t.Run(tc.desc, func(t *testing.T) {
			ret := tc.prompt.dbInsert()
			if ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
			if tc.successfulInsert && tc.prompt.id == 0 {
				t.Fatalf("error. Id for new prompt is not assigned")
			}
		})
	}
	sessionCases := []struct {
		desc     string
		session  session
		expected error
	}{
		{
			"valid session creation",
			session{KsatID: 3},
			nil,
		},
		{
			"invalid session creation (no Ksat)",
			session{KsatID: 10101},
			sql.ErrNoRows,
		},
	}
	for _, tc := range sessionCases {
		t.Run(tc.desc, func(t *testing.T) {
			ret := tc.session.dbInsert()
			if ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
		})
	}
	entryCases := []struct {
		desc     string
		entry    entry
		expected error
	}{
		{
			"valid entry creation",
			entry{sessionID: 2, promptID: 1, txt: "some text"},
			nil,
		},
		{
			"invalid entry creation (no valid sessionID)",
			entry{sessionID: 10101, promptID: 1, txt: "some txt"},
			sql.ErrNoRows,
		},
		{
			"invalid entry creation (no valid promptID)",
			entry{sessionID: 2, promptID: 10101, txt: "some txt"},
			sql.ErrNoRows,
		},
	}
	for _, tc := range entryCases {
		t.Run(tc.desc, func(t *testing.T) {
			ret := tc.entry.dbInsert()
			if ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
		})
	}
}
func TestGetKsats(t *testing.T) {
	type result struct {
		Ksats []Ksat
		err   error
	}
	cases := []struct {
		desc     string
		expected result
	}{
		{
			"gets all cases",
			result{
				Ksats: []Ksat{
					Ksat{id: 1, name: "first", usage: "usage here"},
					Ksat{id: 2, name: "second", usage: "more usage here"},
					Ksat{id: 3, name: "hasP", usage: "usage"},
				},
				err: nil,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			Ksats, err := GetKsats()
			//
			if err != tc.expected.err {
				t.Fatalf("errors don't match: %v, %v", err, tc.expected.err)
			}
			for i, item := range tc.expected.Ksats {
				//  needed since I am creating more sessions than are shown
				if i >= len(Ksats) {
					t.Fatalf("Ksats returned have less items than expected: items missing at ith '%v': %v", i, Ksats[i:])
				}
				if item.id != tc.expected.Ksats[i].id {
					t.Fatalf("ids don't match: %v, %v", item.id, tc.expected.Ksats[i].id)
				}
				if item.name != tc.expected.Ksats[i].name {
					t.Fatalf("name don't match: %v, %v", item.name, tc.expected.Ksats[i].name)
				}
				if item.usage != tc.expected.Ksats[i].usage {
					t.Fatalf("usages don't match: %v, %v", item.usage, tc.expected.Ksats[i].usage)
				}
			}
			// TODO below code will fail. This is due to other tests creating new Ksats and soon to be deleting them
			/* if len(Ksats) > len(tc.expected.Ksats) {
				t.Fatalf("Ksats returned have more items than expected: items missing at ith '%v': %v",
					len(tc.expected.ksats), ksats[len(tc.expected.ksats):])
			}
			*/
		})
	}
}
