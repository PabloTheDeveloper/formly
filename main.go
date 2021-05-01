package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

var db *sql.DB

func createDB() error {
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS ksats (
			ksat_id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			name TEXT UNIQUE NOT NULL,
			usage TEXT NOT NULL,
			CHECK(length(name) >= 1 AND length(name) <= 6 AND length(usage) >= 5 AND length(usage) <= 40)
		);

		CREATE TABLE IF NOT EXISTS prompts (
			prompt_id INTEGER PRIMARY KEY AUTOINCREMENT,
			sequence INTEGER NOT NULL CHECK(sequence >= 0),
			flag TEXT NOT NULL CHECK(length(flag) >= 1 AND length(flag) <= 10),
			usage TEXT NOT NULL CHECK(length(usage) >= 5 AND length(usage) <= 40),
			ksat_id INTEGER NOT NULL,
			FOREIGN KEY (ksat_id) REFERENCES ksats (ksat_id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS sessions (
			session_id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			ksat_id INTEGER NOT NULL,
			FOREIGN KEY (ksat_id) REFERENCES ksats (ksat_id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS entries (
			entry_id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			sequence INTEGER NOT NULL CHECK(sequence > -1),
			label TEXT NOT NULL CHECK(length(label) > 0),
			txt TEXT,
			session_id INTEGER NOT NULL,
			FOREIGN KEY (session_id) REFERENCES sessions (session_id) ON UPDATE CASCADE ON DELETE CASCADE
		);
	`)
	return err
}
func main() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dataPath := filepath.Join(homePath, ".local", "share", "ksat")
	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	dbPath := filepath.Join(dataPath, "data.db")
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := createDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
