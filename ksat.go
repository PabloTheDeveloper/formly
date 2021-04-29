package main

import (
	"database/sql"
	"errors"
	"log"
)

type ksat struct {
	id                int64
	name, description string
}

func getKsatIdByName(name string) (ksatId int64, err error) {
	// TODO add name checking
	err = db.QueryRow("SELECT ksat_id FROM ksats WHERE name = ?", name).Scan(&ksatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	return
}
func newKsat(name, usage, prompts string) error {
	if !(len(name) >= 1 && len(name) <= 6) {
		return errors.New("'name' must be between 1-6 characters long")
	}
	if !(len(usage) >= 5 && len(usage) <= 40) {
		return errors.New("'usage' must be between 5-40 characters")
	}
	ksatId, err := getKsatIdByName(name)
	if err != nil {
		return err
	}
	if ksatId != -1 {
		return errors.New("ksat already exists")
	}
	// TODO parse prompts
	stmt, err := db.Prepare("INSERT INTO ksats (name, usage) VALUES (?, ?)")

	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(name, usage)
	if err != nil {
		log.Fatal(err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
