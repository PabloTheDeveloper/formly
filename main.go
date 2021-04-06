package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

var db *sql.DB

func createDB(clearPrior bool) error {
	if clearPrior {
		_, err := db.Exec(`
			DROP TABLE IF EXISTS ksats;
			DROP TABLE IF EXISTS prompts;
			DROP TABLE IF EXISTS entries;
		`)
		if err != nil {
			return err
		}
	}
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS ksats (
			ksat_id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			name TEXT UNIQUE NOT NULL,
			description TEXT NOT NULL,
			CHECK(length(name) > 1 AND length(name) < 7 AND length(description) > 5)
		);

		CREATE TABLE IF NOT EXISTS prompts (
			prompt_id INTEGER PRIMARY KEY AUTOINCREMENT,
			sequence INTEGER NOT NULL CHECK(sequence > -1),
			label TEXT NOT NULL CHECK(length(label) > 0 AND length(label) < 7),
			input_type TEXT CHECK(input_type IN ('str','num','flt','time', 'audio', 'video', 'path', 'binary')) NOT NULL DEFAULT 'str',
			ksat_id INTEGER NOT NULL,
			FOREIGN KEY (ksat_id) REFERENCES ksats (ksat_id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS schedules (
			schedule_id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			once DATETIME,
			selected INTEGER NOT NULL,
			ksat_id INTEGER NOT NULL,
			FOREIGN KEY (ksat_id) REFERENCES ksats (ksat_id)
		);

		CREATE TABLE IF NOT EXISTS recurrings (
			recurring_id INTEGER PRIMARY KEY AUTOINCREMENT,
			day TEXT CHECK(day in ('mon', 'tue', 'wed', 'thu', 'fri', 'sat', 'sun')) DEFAULT 'sun',
			hour INTEGER CHECK(hour < 25 AND hour > -2),
			min INTEGER CHECK(min < 25 AND min > -2)
		);

		CREATE TABLE IF NOT EXISTS schedules_recurrings (
			schedule_id INTEGER NOT NULL,
			recurring_id INTEGER NOT NULL,
			FOREIGN KEY (schedule_id) REFERENCES schedules (schedule_id),
			FOREIGN KEY (recurring_id) REFERENCES recurrings (recurring_id)
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
			input_type TEXT CHECK(input_type IN ('str','num','flt','time', 'audio', 'video', 'path', 'binary')) NOT NULL DEFAULT 'str',
			txt TEXT,
			int INTEGER,
			flt FLOAT,
			bin BLOB,
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
	if err := createDB(true); err != nil {
		log.Fatal(err)
	}
	if err := execute(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}
