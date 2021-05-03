package main

import (
	"testing"
)

/*** Unit Tests ***/
func TestValidateName(t *testing.T) {
	cases := []struct {
		desc     string
		name     string
		expected error
	}{
		{
			"ZeroLetterWord",
			"",
			strLengthErr{lower: 1, upper: 6, str: ""},
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
			"SixLetterWord",
			"sixsix",
			nil,
		},
		{
			"SevenLetterWord",
			"sevsevs",
			strLengthErr{lower: 1, upper: 6, str: "sevsevs"},
		},
		{
			"WordWithNumberAndLetters",
			"a12b",
			wordErr{word: "a12b"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			task := ksat{name: tc.name}
			if ret := task.validateName(); ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
		})
	}
}
func TestValidateUsage(t *testing.T) {
	cases := []struct {
		desc     string
		usage    string
		expected error
	}{
		{
			"ZeroLetterUsage",
			"",
			strLengthErr{lower: 5, upper: 40, str: ""},
		},
		{
			"FourLetterUsage",
			"abcd",
			strLengthErr{lower: 5, upper: 40, str: "abcd"},
		},
		{
			"FiveLetterUsage",
			"abcde",
			nil,
		},
		{
			"SixLetterUsage",
			"sixsix",
			nil,
		},
		{
			"FourtyLetterUsage",
			"0123456789" +
				"0123456789" +
				"0123456789" +
				"0123456789",
			nil,
		},
		{
			"FourtyOneLetterUsage",
			"0123456789" +
				"0123456789" +
				"0123456789" +
				"0123456789" + "1",
			strLengthErr{lower: 5, upper: 40,
				str: "0123456789" +
					"0123456789" +
					"0123456789" +
					"0123456789" + "1"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			task := ksat{usage: tc.usage}
			if ret := task.validateUsage(); ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
		})
	}
}
func TestValidate(t *testing.T) {
	cases := []struct {
		desc     string
		task     ksat
		success  bool
		expected error
	}{
		{
			"A Valid ksat",
			ksat{name: "aName", usage: "0123456789" +
				"0123456789" +
				"0123456789" +
				"0123456789"},
			true,
			nil,
		},
		{
			"A Valid ksat",
			ksat{name: "", usage: "1234"},
			false,
			strLengthErr{lower: 1, upper: 6, str: ""},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			if ret := tc.task.validate(); ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
		})
	}
}

/*** Integration Tests ***/
func TestGetKsatByName(t *testing.T) {
	cases := []struct {
		desc       string
		task       ksat
		successful ksat
		expected   error
	}{
		{
			"valid name but for a ksat that exists",
			ksat{name: "first", usage: "first usage here"},
			ksat{id: 1, name: "first", usage: "some usage"},
			nil,
		},
		{
			"another valid name but for a ksat that exists",
			ksat{name: "second", usage: "second usage here"},
			ksat{id: 2, name: "second", usage: "some more usage"},
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			ret := tc.task.getKsatByName()
			if ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
			if tc.task.id != tc.successful.id {
				t.Fatalf("ids don't match: %v, %v", tc.task.id, tc.successful.id)
			}
			if tc.task.name != tc.successful.name {
				t.Fatalf("names don't match: %v, %v", tc.task.name, tc.successful.name)
			}
			if tc.task.usage != tc.successful.usage {
				t.Fatalf("usages don't match: %v, %v", tc.task.usage, tc.successful.usage)
			}
		})
	}
}

func TestDbInsert(t *testing.T) {
	cases := []struct {
		desc             string
		task             ksat
		successfulInsert bool
		expected         error
	}{
		{
			"invalid name for ksat (letter + number)",
			ksat{name: "a12b", usage: "some usage"},
			false,
			wordErr{word: "a12b"},
		},
		{
			"valid name but for a ksat that exists",
			ksat{name: "first", usage: "first usage here"},
			false,
			alreadyExistsErr{identifier: "first", tableName: "ksat"},
		},
		{
			"valid name but for a ksat that does not exists",
			ksat{name: "new", usage: "usage here"},
			true,
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			ret := tc.task.dbInsert()
			if ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
			if tc.successfulInsert && tc.task.id == 0 {
				t.Fatalf("error. Id for new ksat is not assigned")
			}
		})
	}
}
