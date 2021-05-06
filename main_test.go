package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func setUp() (string, error) {
	dir, err := ioutil.TempDir("", "testDir-*")
	if err != nil {
		return dir, err
	}
	db, err = sql.Open("sqlite3", filepath.Join(dir, "data.db"))
	if err != nil {
		return dir, err
	}
	if err := createDB(); err != nil {
		return dir, err
	}
	/*
		Creating ksats in db for testing
	*/
	createKsatStmt, err := db.Prepare("INSERT INTO ksats (name, usage) VALUES (?, ?)")
	if err != nil {
		return dir, err
	}
	ksats := []ksat{
		ksat{
			name:  "first",
			usage: "some usage",
		},
		ksat{
			name:  "second",
			usage: "some more usage",
		},
		ksat{
			name:  "hasP",
			usage: "usage",
		},
	}
	for _, task := range ksats {
		if _, err := createKsatStmt.Exec(task.name, task.usage); err != nil {
			return dir, err
		}
	}
	/*
		Creating prompts in db for testing
	*/
	createPromptStmt, err := db.Prepare("INSERT INTO prompts (ksat_id, sequence, flag, usage) VALUES (?, ?, ?, ?)")
	if err != nil {
		return dir, err
	}
	prompts := []prompt{
		prompt{
			ksatID:   3,
			sequence: 1,
			flag:     "firstflag",
			usage:    "some usage",
		},
		prompt{
			ksatID:   3,
			sequence: 2,
			flag:     "secondflag",
			usage:    "some usage",
		},
	}
	for _, item := range prompts {
		if _, err := createPromptStmt.Exec(item.ksatID, item.sequence, item.flag, item.usage); err != nil {
			return dir, err
		}
	}
	/*
		Creating sessions in db for testing
	*/
	createSessionStmt, err := db.Prepare("INSERT INTO sessions (ksat_id) VALUES (?)")
	if err != nil {
		return dir, err
	}
	sessions := []session{{ksatID: 3}, {ksatID: 3}}
	for _, session := range sessions {
		if _, err := createSessionStmt.Exec(session.ksatID); err != nil {
			return dir, err
		}
	}
	return dir, nil
}

func tearDown(dir string) error {
	if err := db.Close(); err != nil {
		return err
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	tmpDir, err := setUp()
	if err != nil {
		log.Fatal(err)
	}
	exitVal := m.Run()
	if err := tearDown(tmpDir); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitVal)
}
