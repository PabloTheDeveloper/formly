package main

import (
	"fmt"
	"math"
)

type prompt struct {
	id, ksatID, sequence int64
	flag, usage          string
}

func (prompt *prompt) validateSequence() error {
	return isNumBoundsCorrect(prompt.sequence, 1, math.MaxInt64)
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

type flagNotUniqueErr struct {
	flag string
}

func (e flagNotUniqueErr) Error() string {
	return fmt.Sprintf("'%s' is not unique", e.flag)
}
func validatePromptsFlags(flag string, prompts []prompt) error {
	p := prompt{flag: flag}
	if err := p.validateFlag(); err != nil {
		return err
	}
	for _, item := range prompts {
		if flag == item.flag {
			return flagNotUniqueErr{flag: flag}
		}
	}
	return nil
}

type sequenceErr struct {
	sequence, accurateSequence int64
	prompt                     prompt
}

func (e sequenceErr) Error() string {
	return fmt.Sprintf("'%v' is the sequence for prompt: %v.\nIt should be %v",
		e.sequence, e.prompt, e.accurateSequence)
}
func validatePromptsSequences(prompts []prompt) error {
	for i, item := range prompts {
		if expectedSequence := int64(i + 1); expectedSequence != item.sequence {
			return sequenceErr{sequence: item.sequence, accurateSequence: expectedSequence, prompt: item}
		}
	}
	return nil
}
func (prompt *prompt) getByID() error {
	return db.QueryRow(
		"SELECT prompt_id, ksat_id, sequence, flag, usage FROM prompts WHERE prompt_id = ?",
		prompt.id,
	).Scan(
		&prompt.id,
		&prompt.ksatID,
		&prompt.sequence,
		&prompt.flag,
		&prompt.usage,
	)
}
func (prompt *prompt) dbInsert() error {
	if err := prompt.validate(); err != nil {
		return err
	}
	task := ksat{id: prompt.ksatID}
	prompts, err := task.getPromptsByID()
	if err != nil {
		return err
	}
	if err := validatePromptsFlags(prompt.flag, prompts); err != nil {
		return err
	}
	if err := validatePromptsSequences(prompts); err != nil {
		return err
	}
	if expectedSequence := int64(len(prompts) + 1); expectedSequence != prompt.sequence {
		return sequenceErr{sequence: prompt.sequence, accurateSequence: expectedSequence, prompt: *prompt}
	}
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
