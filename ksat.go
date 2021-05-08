package main

import (
	"database/sql"
	"fmt"
)

type Ksat struct {
	id          int64
	name, usage string
}

func (ksat *Ksat) SetName(name string) error {
	if err := isStringLengthCorrect(name, 1, 6); err != nil {
		return err
	}
	if err := isWordValid(name); err != nil {
		return err
	}
	ksat.name = name
	return nil
}
func (ksat *Ksat) SetUsage(usage string) error {
	return isStringLengthCorrect(usage, 5, 40)
}
func GetKsatByName(name string) (ksat Ksat, err error) {
	if err = ksat.SetName(name); err != nil {
		return Ksat{}, err
	}
	err = db.QueryRow(
		"select ksat_id, name, usage from ksats where name = ?", ksat.name,
	).Scan(&ksat.id, &ksat.name, &ksat.usage)
	return
}

func GetKsatByID(ksatID int64) (ksat Ksat, err error) {
	err = db.QueryRow(
		"select ksat_id, name, usage from ksats where ksat_id = ?",
		ksatID,
	).Scan(&ksat.id, &ksat.name, &ksat.usage)
	return
}
func GetPromptsByKsatID(ksatID int64) ([]prompt, error) {
	if _, err := GetKsatByID(ksatID); err != nil {
		return nil, err
	}
	prompts := []prompt{}
	rows, err := db.Query("select ksat_id, prompt_id, sequence, flag, usage from prompts where ksat_id = ?", ksatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := prompt{}
		err := rows.Scan(&item.KsatID, &item.id, &item.sequence, &item.flag, &item.usage)
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
func GetSessionsByKsatID(ksatID int64) ([]session, error) {
	if _, err := GetKsatByID(ksatID); err != nil {
		return nil, err
	}
	sessions := []session{}
	rows, err := db.Query("SELECT session_id, Ksat_id, created_at FROM sessions WHERE Ksat_id = ?", ksatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := session{}
		err := rows.Scan(&item.id, &item.KsatID, &item.createAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return sessions, nil
}
func GetKsats() ([]Ksat, error) {
	Ksats := []Ksat{}
	rows, err := db.Query("SELECT Ksat_id, name, usage FROM Ksats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := Ksat{}
		err := rows.Scan(&item.id, &item.name, &item.usage)
		if err != nil {
			return nil, err
		}
		Ksats = append(Ksats, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return Ksats, nil
}

type AlreadyExistsErr struct {
	identifier, tableName string
}

func (e AlreadyExistsErr) Error() string {
	return fmt.Sprintf("'%s', %s Already exists", e.identifier, e.tableName)
}

// assumes no setting of fields without appropriate setters
func (ksat *Ksat) dbInsert() error {
	if ksat, err := GetKsatByName(ksat.name); err == nil {
		return AlreadyExistsErr{identifier: ksat.name, tableName: "Ksat"}
	} else if err != sql.ErrNoRows {
		return err // some other error than Ksat error
	}
	stmt, err := db.Prepare("INSERT INTO Ksats (name, usage) VALUES (?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(ksat.name, ksat.usage)
	if err != nil {
		return err
	}
	ksat.id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
