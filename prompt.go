package main

import "math"

type prompt struct {
	id, ksatID, sequence int64
	flag, usage          string
}

func (prompt *prompt) validateSequence() error {
	return isNumBoundsCorrect(prompt.sequence, 0, math.MaxInt64)
}
func (prompt *prompt) validateFlag() error {
	if err := isStringLengthCorrect(prompt.flag, 1, 10); err != nil {
		return err
	}
	return isWordValid(prompt.flag)
}
func (prompt *prompt) validateUsage() error {
	return isStringLengthCorrect(prompt.usage, 5, 40)
}
func (prompt *prompt) validate() error {
	if err := prompt.validateSequence(); err != nil {
		return err
	}
	if err := prompt.validateFlag(); err != nil {
		return err
	}
	return prompt.validateUsage()
}
func (prompt *prompt) dbInsert() error {
	if err := prompt.validate(); err != nil {
		return err
	}
	task := ksat{id: prompt.ksatID}
	if err := task.getByID(); err != nil {
		return err
	}
	//TODO ensure that no other prompt with ksat contains sequence
	// TODO ensure flag with ksat combination is unique
	stmt, err := db.Prepare("INSERT INTO prompts (ksat_id, sequence, flag, usage) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(prompt.ksatID, prompt.sequence, prompt.flag, prompt.usage)
	if err != nil {
		return err
	}
	prompt.id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
