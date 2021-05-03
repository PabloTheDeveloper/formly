package main

import (
	"fmt"
	"regexp"
)

type wordErr struct {
	word string
}

func (e wordErr) Error() string {
	return fmt.Sprintf("'%v' word must contain only letters", e.word)
}
func isWordValid(word string) error {
	if regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(word) {
		return nil
	}
	return wordErr{word: word}
}

type strLengthErr struct {
	lower, upper int
	str          string
}

func (e strLengthErr) Error() string {
	return fmt.Sprintf("'%s' must be between %v - %v characters long", e.str, e.lower, e.upper)
}

func isStringLengthCorrect(str string, lower, upper int) error {
	if len(str) < lower || len(str) > upper {
		return strLengthErr{lower: lower, upper: upper, str: str}
	}
	return nil
}
