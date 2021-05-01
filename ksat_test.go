package main

import "testing"

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
func TestIsKsatNameValid(t *testing.T) {
	cases := []struct {
		desc     string
		name     string
		expected error
	}{
		{
			"ZeroLetterWord",
			"",
			ksatNameErr{name: ""},
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
			ksatNameErr{name: "sevsevs"},
		},
		{
			"WordWithNumberAndLetters",
			"a12b",
			wordErr{word: "a12b"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			if ret := isKsatNameValid(tc.name); ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
		})
	}
}
func TestIsKsatUsageValid(t *testing.T) {
	cases := []struct {
		desc     string
		usage    string
		expected error
	}{
		{
			"ZeroLetterUsage",
			"",
			ksatUsageErr{usage: ""},
		},
		{
			"FourLetterUsage",
			"abcd",
			ksatUsageErr{usage: "abcd"},
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
			ksatUsageErr{
				usage: "0123456789" +
					"0123456789" +
					"0123456789" +
					"0123456789" + "1"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			if ret := isKsatUsageValid(tc.usage); ret != tc.expected {
				t.Fatalf("errors don't match: %v, %v", ret, tc.expected)
			}
		})
	}
}
func TestValidate(t *testing.T) {
	cases := []struct {
		desc     string
		task     ksat
		expected error
	}{
		{
			"ZeroLetterInvalidNameValidUsage",
			ksat{name: "", usage: "some usage"},
			ksatNameErr{name: ""},
		},
		{
			"ZeroLetterInvalidUsageValidName",
			ksat{name: "aName", usage: ""},
			ksatUsageErr{usage: ""},
		},
		{
			"FourLetterInvalidUsageValidName",
			ksat{name: "aName", usage: "abcd"},
			ksatUsageErr{usage: "abcd"},
		},
		{
			"SixLetterValidUsageSixLetterValidName",
			ksat{name: "sixsix", usage: "sixsix"},
			nil,
		},
		{
			"FourtyLetterUsageValidName",
			ksat{name: "aName", usage: "0123456789" +
				"0123456789" +
				"0123456789" +
				"0123456789"},
			nil,
		},
		{
			"FourtyOneLetterUsage",
			ksat{name: "aName", usage: "0123456789" +
				"0123456789" +
				"0123456789" +
				"0123456789" + "1"},
			ksatUsageErr{
				usage: "0123456789" +
					"0123456789" +
					"0123456789" +
					"0123456789" + "1"},
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
