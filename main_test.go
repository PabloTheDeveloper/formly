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
		return "", err
	}
	db, err = sql.Open("sqlite3", filepath.Join(dir, "data.db"))
	if err != nil {
		return "", err
	}
	err = createDB(true)
	return dir, err
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
