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
	}
	for _, task := range ksats {
		if _, err := createKsatStmt.Exec(task.name, task.usage); err != nil {
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
