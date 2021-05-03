package main

import (
	"fmt"
	"regexp"
)

type UserError interface {
	error
	UserError() string
}

type wordErr struct {
	word string
}

func (e wordErr) Error() string {
	return fmt.Sprintf("'%v' word must contain only letters", e.word)
}
func (e wordErr) UserError() string {
	return e.Error()
}
func isWordValid(word string) error {
	if regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(word) {
		return nil
	}
	return wordErr{word: word}
}

type ksatNameErr struct {
	name string
}

func (e ksatNameErr) Error() string {
	return fmt.Sprintf("'%v' must be between 1-6 characters long", e.name)
}
func (e ksatNameErr) UserError() string {
	return e.Error()
}
func isKsatNameValid(name string) error {
	if !(len(name) >= 1 && len(name) <= 6) {
		return ksatNameErr{name: name}
	}
	if err := isWordValid(name); err != nil {
		return err
	}
	return nil
}

type ksatUsageErr struct {
	usage string
}

func (e ksatUsageErr) Error() string {
	return fmt.Sprintf("'%v' must be between 5-40 characters long", e.usage)
}
func (e ksatUsageErr) UserError() string {
	return e.Error()
}
func isKsatUsageValid(usage string) error {
	if !(len(usage) >= 5 && len(usage) <= 40) {
		return ksatUsageErr{usage: usage}
	}
	return nil
}

type ksat struct {
	id          int64
	name, usage string
	isValidated bool
}

func (task *ksat) validate() error {
	if err := isKsatNameValid(task.name); err != nil {
		task.isValidated = false
		return err
	}
	if err := isKsatUsageValid(task.usage); err != nil {
		task.isValidated = false
		return err
	}
	task.isValidated = true
	return nil
}

type ksatDbInsertErr struct {
	name string
}

func (e ksatDbInsertErr) Error() string {
	return fmt.Sprintf("'%v' already exists. Please select a new name", e.name)
}
func (e ksatDbInsertErr) UserError() string {
	return e.Error()
}

func (task *ksat) dbInsert() error {
	if err := task.validate(); err != nil {
		return err
	}
	// ensure name does not already exists
	if id, err := getKsatIdByName(task.name); id != 0 {
		return ksatDbInsertErr{name: task.name}
	} else if _, ok := err.(noKsatIdByNameErr); !ok {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO ksats (name, usage) VALUES (?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(task.name, task.usage)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	task.id = id
	return nil
}
