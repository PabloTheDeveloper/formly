package main

import (
	"database/sql"
)

// getKsatIdByName if returns -1 means it was not found or error along way
func getKsatIdByName(name string) (int64, error) {
	var id int64 = -1
	if err := isKsatNameValid(name); err != nil {
		return id, err
	}
	err := db.QueryRow("SELECT ksat_id FROM ksats WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		err = nil
	}
	return id, err
}

/*
func (task *ksat) dbCreate() error {
	stmt, err := db.Prepare("INSERT INTO ksats (name, usage) VALUES (?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(task.name, task.usage)
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
*/
