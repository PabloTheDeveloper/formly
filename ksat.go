package main

import (
	"database/sql"
	"fmt"
)

type ksat struct {
	id          int64
	name, usage string
	isValidated bool
}

func (task *ksat) validateName() error {
	if err := isStringLengthCorrect(task.name, 1, 6); err != nil {
		return err
	}
	return isWordValid(task.name)
}
func (task *ksat) validateUsage() error {
	return isStringLengthCorrect(task.usage, 5, 40)
}
func (task *ksat) validate() error {
	task.isValidated = false
	if err := task.validateName(); err != nil {
		return err
	}
	if err := task.validateUsage(); err != nil {
		return err
	}
	task.isValidated = true
	return nil
}
func (task *ksat) getByName() error {
	if err := task.validateName(); err != nil {
		return err
	}
	return db.QueryRow(
		"select ksat_id, name, usage from ksats where name = ?", task.name,
	).Scan(&task.id, &task.name, &task.usage)
}
func (task *ksat) getByID() error {
	return db.QueryRow(
		"select ksat_id, name, usage from ksats where ksat_id = ?", task.id).Scan(&task.id, &task.name, &task.usage)
}

type alreadyExistsErr struct {
	identifier, tableName string
}

func (e alreadyExistsErr) Error() string {
	return fmt.Sprintf("'%s', %s already exists", e.identifier, e.tableName)
}

func (task *ksat) getPromptsByID() ([]prompt, error) {
	if err := task.getByID(); err != nil {
		return nil, err
	}
	prompts := []prompt{}
	rows, err := db.Query("select ksat_id, prompt_id, sequence, flag, usage from prompts where ksat_id = ?", task.id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := prompt{}
		err := rows.Scan(&item.ksatID, &item.id, &item.sequence, &item.flag, &item.usage)
		if err != nil {
			return nil, err
		}
		prompts = append(prompts, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return prompts, nil
}
func (task *ksat) dbInsert() error {
	if err := task.validate(); err != nil {
		return err
	}
	if err := task.getByName(); err == nil {
		return alreadyExistsErr{identifier: task.name, tableName: "ksat"}
	} else if err != sql.ErrNoRows {
		return err // some other error than ksat error
	}
	stmt, err := db.Prepare("INSERT INTO ksats (name, usage) VALUES (?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(task.name, task.usage)
	if err != nil {
		return err
	}
	task.id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
func getKsats() ([]ksat, error) {
	ksats := []ksat{}
	rows, err := db.Query("SELECT ksat_id, name, usage FROM ksats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := ksat{}
		err := rows.Scan(&item.id, &item.name, &item.usage)
		if err != nil {
			return nil, err
		}
		ksats = append(ksats, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ksats, nil
}
