package main

import (
	"database/sql"
	"errors"
	"log"
)

type ksat struct {
	id                int64
	name, description string
	schedule          schedule
	prompts           []prompt
	sessions          []session
}
type prompt struct {
	id               int64
	label, inputType string
}
type schedule struct {
	id         int64
	once       string
	selected   bool
	reccurings []reccuring
}
type reccuring struct {
	day  int64
	hour int
	min  int
}
type session struct {
	id        int64
	createdAt string
	entries   []entry
}
type entry struct {
	id               int64
	label, inputType string
	txt              string
	integer          int
	flt              float64
	bin              []byte
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
func newKsat(name, description string) (ksat, error) {
	if len(name) <= 1 || len(name) >= 7 {
		return ksat{}, errors.New("'name' must be between 1-6 characters long")
	}
	if len(description) <= 5 {
		return ksat{}, errors.New("'desc' must be larger than 5 characters long")
	}
	stmt, err := db.Prepare("INSERT INTO ksats (name, description) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(name, description)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return ksat{lastId, name, description, schedule{}, nil, nil}, nil
}
