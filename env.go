package ksat

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	// to support sqlite
	_ "github.com/mattn/go-sqlite3"
)

const (
	localStorage uint8 = 1 << iota
	sqlite
)

type storageOption uint8

// LocalSqlite ...
const LocalSqlite storageOption = storageOption(localStorage | sqlite)

// Env ...
type Env struct {
	FormModel
	LabelModel
	SubmissionModel
	EntryModel
	close func() error
}

// FormModel ...
type FormModel interface {
	GetByName(name string) (Form, error)
	GetByID(id int64) (Form, error)
	GetAll() ([]Form, error)
}

// LabelModel ...
type LabelModel interface {
	Create(formID, position int64, name, usage string) (Label, error)
	GetLabels(formID int64) ([]Label, error)
}

// SubmissionModel ...
type SubmissionModel interface {
	Create(formID int64) (Submission, error)
}

// EntryModel ...
type EntryModel interface {
	Create(submissionID, labelID int64, txt string) (Entry, error)
}

// NewEnv ...
func NewEnv(op storageOption) (*Env, error) {
	switch op {
	case LocalSqlite:
		db, err := newLocalSqLite()
		if err != nil {
			return nil, err
		}
		return newSQLEnv(db), nil
	default:
		return nil, fmt.Errorf("unsupported storage option")
	}
}

// Close ...
func (env *Env) Close() error {
	return env.close()
}
func newLocalSqLite() (*sql.DB, error) {
	schema := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE IF NOT EXISTS forms (
			form_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL CHECK(length(name) >= 1 AND length(name) <= 6),
			usage TEXT NOT NULL CHECK(length(usage) >= 5 AND length(usage) <= 40),
			editable BOOL DEFAULT TRUE,
			deletable BOOL DEFAULT TRUE
		);

		CREATE TABLE IF NOT EXISTS labels (
			label_id INTEGER PRIMARY KEY AUTOINCREMENT,
			position INTEGER NOT NULL CHECK(position >= 1),
			name TEXT NOT NULL CHECK(length(name) >= 1 AND length(name) <= 10),
			usage TEXT NOT NULL CHECK(length(usage) >= 5 AND length(usage) <= 40),
			repeatable BOOL DEFAULT FALSE,
			editable BOOL DEFAULT TRUE,
			deletable BOOL DEFAULT TRUE,
			form_id INTEGER NOT NULL,
			FOREIGN KEY (form_id) REFERENCES forms (form_id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS submissions (
			submission_id INTEGER PRIMARY KEY AUTOINCREMENT,
			form_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (form_id) REFERENCES forms (form_id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS entries (
			entry_id INTEGER PRIMARY KEY AUTOINCREMENT,
			label_id INTEGER NOT NULL,
			submission_id INTEGER NOT NULL,
			txt TEXT,
			FOREIGN KEY (label_id) REFERENCES labels (label_id) ON UPDATE CASCADE ON DELETE CASCADE,
			FOREIGN KEY (submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		INSERT INTO forms(name, usage, editable, deletable)
		SELECT 'create', 'subcommand to create other tasks', FALSE, FALSE
		WHERE NOT EXISTS(SELECT 1 FROM forms WHERE name = 'create');

			INSERT INTO labels(name, usage, editable, deletable,position, form_id)
			SELECT 'name', 'name for new form', FALSE, FALSE, 1, 1
			WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 1);

			INSERT INTO labels(name, usage, editable, deletable,position, form_id)
			SELECT 'labels', 'comma seperate list of labels for form', FALSE, FALSE, 2, 1
			WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 2);

		INSERT INTO forms(name, usage, editable, deletable)
		SELECT 'update', 'subcommand to update tasks', FALSE, FALSE
		WHERE NOT EXISTS(SELECT 1 FROM forms WHERE name = 'update');

			INSERT INTO labels(name, usage, editable, deletable,position, form_id)
			SELECT 'name', 'name for new form', FALSE, FALSE, 1, 2
			WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 3);

			INSERT INTO labels(name, usage, editable, deletable,position, form_id)
			SELECT 'labels', 'comma seperate list of labels for form', FALSE, FALSE, 2, 2
			WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 4);

		INSERT INTO forms(name, usage, editable, deletable)
		SELECT 'delete', 'subcommand to delete other tasks', FALSE, FALSE
		WHERE NOT EXISTS(SELECT 1 FROM forms WHERE name = 'delete');

			INSERT INTO labels(name, usage, editable, deletable,position, form_id)
			SELECT 'name', 'name for new form', FALSE, FALSE, 1, 3
			WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 5);

			INSERT INTO labels(name, usage, editable, deletable,position, form_id)
			SELECT 'labels', 'comma seperate list of labels for form', FALSE, FALSE, 2, 3
			WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 6);

		INSERT INTO forms(name, usage, editable, deletable)
		SELECT 'read', 'subcommand to read other tasks', FALSE, FALSE
		WHERE NOT EXISTS(SELECT 1 FROM forms WHERE name = 'read');

			INSERT INTO labels(name, usage, editable, deletable,position, form_id)
			SELECT 'name', 'name for new form', FALSE, FALSE, 1, 4
			WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 7);

			INSERT INTO labels(name, usage, editable, deletable,position, form_id)
			SELECT 'labels', 'comma seperate list of labels for form', FALSE, FALSE, 2, 4
			WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 8);

		`
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dataPath := filepath.Join(homePath, ".local", "share", "ksat")
	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(dataPath, "data.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return db, nil
}

// Form ...
type Form struct {
	id          int64
	name, usage string
}

// GetID gets id from Form
func (form Form) GetID() int64 {
	return form.id
}

// GetName gets name from Form
func (form Form) GetName() string {
	return form.name
}

// GetUsage gets usage from Form
func (form Form) GetUsage() string {
	return form.usage
}

type Label struct {
	id, formID, position int64
	name, usage          string
}

// GetID gets id from Label
func (label Label) GetID() int64 {
	return label.id
}

// GetName gets name from Label
func (label Label) GetName() string {
	return label.name
}

// GetUsage gets usage from Label
func (label Label) GetUsage() string {
	return label.usage
}

type Submission struct {
	id, formID int64
	createAt   time.Time
}

// GetID gets id from Submission
func (submission Submission) GetID() int64 {
	return submission.id
}

type Entry struct {
	id, labelID, submissionID int64
	txt                       string
}

// GetID gets id from Entry
func (entry Entry) GetID() int64 {
	return entry.id
}
