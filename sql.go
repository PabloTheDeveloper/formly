package ksat

import (
	"database/sql"
	"os"
	"path/filepath"
)

type sqlModels struct {
	db *sql.DB
	sqlFormModel
	sqlSubmissionModel
	sqlEntryModel
}

// NewLocalSqLiteEnv ...
func NewLocalSqLiteEnv() (*Env, error) {
	schema := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE IF NOT EXISTS forms (
			form_id INTEGER PRIMARY KEY AUTOINCREMENT,
			editable BOOL DEFAULT TRUE,
			deleteable BOOL DEFAULT TRUE,
			name TEXT NOT NULL CHECK(length(name) >= 1 AND length(name) <= 16),
			usage TEXT NOT NULL CHECK(length(usage) >= 5 AND length(usage) <= 252)
		);

		CREATE TABLE IF NOT EXISTS labels (
			label_id INTEGER PRIMARY KEY AUTOINCREMENT,
			form_id INTEGER NOT NULL,
			position INTEGER NOT NULL CHECK(position >= 1),
			repeatable BOOL DEFAULT FALSE,
			editable BOOL DEFAULT TRUE,
			deleteable BOOL DEFAULT TRUE,
			name TEXT NOT NULL CHECK(length(name) >= 1 AND length(name) <= 16),
			usage TEXT NOT NULL CHECK(length(usage) >= 5 AND length(usage) <= 252),
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
			submission_id INTEGER NOT NULL,
			label_id INTEGER NOT NULL,
			txt TEXT,
			FOREIGN KEY (label_id) REFERENCES labels (label_id) ON UPDATE CASCADE ON DELETE CASCADE,
			FOREIGN KEY (submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		INSERT INTO forms(editable, deleteable, name, usage)
		SELECT FALSE, FALSE, 'create', 'subcommand to create other forms'
		WHERE NOT EXISTS(SELECT 1 FROM forms WHERE name = 'create');

		INSERT INTO labels(form_id, position, editable, deleteable, name, usage)
		SELECT 1, 1, FALSE, FALSE, 'name', 'what the new form name will be'
		WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 1);

		INSERT INTO labels(form_id, position, editable, deleteable, name, usage)
		SELECT 1, 2, FALSE, FALSE, 'usage', 'what the new form usage will be'
		WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 2);

		INSERT INTO labels(form_id, position, editable, deleteable, name, usage)
		SELECT 1, 3, FALSE, FALSE, 'labels', ' what the new labels will be. requiring this str format: [{name:newName, usage:newUsage, repeatable:1_Or_0_default_is_0}, {...}, ...]'
		WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 3);

		INSERT INTO forms(editable, deleteable, name, usage)
		SELECT FALSE, FALSE, 'read', 'subcommand to read form entries'
		WHERE NOT EXISTS(SELECT 1 FROM forms WHERE name = 'read');

		INSERT INTO labels(form_id, position, editable, deleteable, name, usage)
		SELECT 2, 1, FALSE, FALSE, 'name', 'what the new form name will be'
		WHERE NOT EXISTS(SELECT 1 FROM labels WHERE label_id = 4);
		`
	/*


	 */
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
	return &Env{
		FormModel:       sqlFormModel{db: db},
		LabelModel:      sqlLabelModel{db: db},
		SubmissionModel: sqlSubmissionModel{db: db},
		EntryModel:      sqlEntryModel{db: db},
		close: func() error {
			return db.Close()
		},
	}, nil
}

type sqlFormModel struct {
	db *sql.DB
}

func (model sqlFormModel) Create(name, usage string) (Form, error) {
	form := Form{name: name, usage: usage}
	stmt, err := model.db.Prepare(
		"INSERT INTO forms (name, usage) VALUES(?, ?)",
	)
	if err != nil {
		return Form{}, err
	}
	res, err := stmt.Exec(form.name, form.usage)
	if err != nil {
		return Form{}, err
	}
	form.id, err = res.LastInsertId()
	if err != nil {
		return Form{}, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return Form{}, err
	}
	return form, nil
}
func (model sqlFormModel) GetByName(name string) (Form, error) {
	// return len(name) >= 1 && len(name) <= 6 && regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(name)
	form := Form{}
	err := model.db.QueryRow(
		"SELECT form_id, name, usage FROM forms WHERE name = ?",
		name,
	).Scan(&form.id, &form.name, &form.usage)
	if err == sql.ErrNoRows {
		err = nil
	}
	return form, err
}
func (model sqlFormModel) GetByID(id int64) (Form, error) {
	form := Form{}
	err := model.db.QueryRow(
		"SELECT form_id, name, usage FROM forms WHERE id = ?",
		id,
	).Scan(&form.id, &form.name, &form.usage)
	if err == sql.ErrNoRows {
		err = nil
	}
	return form, err
}
func (model sqlFormModel) GetAll() ([]Form, error) {
	forms := []Form{}
	rows, err := model.db.Query("SELECT form_id, name, usage FROM forms")
	if err != nil {
		return forms, err
	}
	defer rows.Close()
	for rows.Next() {
		form := Form{}
		err := rows.Scan(&form.id, &form.name, &form.usage)
		if err != nil {
			return forms, err
		}
		forms = append(forms, form)
	}
	err = rows.Err()
	return forms, err
}

type sqlLabelModel struct {
	db *sql.DB
}

func (model sqlLabelModel) Create(formID, position int64, repeatable bool, name, usage string) (Label, error) {
	label := Label{formID: formID, position: position, repeatable: repeatable, name: name, usage: usage}
	stmt, err := model.db.Prepare(
		"INSERT INTO labels (form_id, position, repeatable, name, usage) VALUES(?, ?, ?, ?, ?)",
	)
	if err != nil {
		return Label{}, err
	}
	res, err := stmt.Exec(label.formID, label.position, label.repeatable, label.name, label.usage)
	if err != nil {
		return Label{}, err
	}
	label.id, err = res.LastInsertId()
	if err != nil {
		return Label{}, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return Label{}, err
	}
	return label, nil
}

func (model sqlLabelModel) GetLabels(formID int64) ([]Label, error) {
	labels := []Label{}
	rows, err := model.db.Query("SELECT label_id, position, repeatable, name, usage FROM labels WHERE form_id = ? ORDER BY position ASC", formID)
	if err != nil {
		return labels, err
	}
	defer rows.Close()
	for rows.Next() {
		label := Label{}
		err := rows.Scan(&label.id, &label.position, &label.repeatable, &label.name, &label.usage)
		if err != nil {
			return labels, err
		}
		labels = append(labels, label)
	}
	err = rows.Err()
	return labels, err
}

type sqlSubmissionModel struct {
	db *sql.DB
}

func (model sqlSubmissionModel) Create(formID int64) (Submission, error) {
	submission := Submission{formID: formID}
	stmt, err := model.db.Prepare("INSERT INTO submissions (form_id) VALUES(?)")
	if err != nil {
		return Submission{}, err
	}
	res, err := stmt.Exec(submission.formID)
	if err != nil {
		return Submission{}, err
	}
	submission.id, err = res.LastInsertId()
	if err != nil {
		return Submission{}, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return Submission{}, err
	}
	return submission, nil
}
func (model sqlSubmissionModel) GetSubmissions(formID int64) ([]Submission, error) {
	submissions := []Submission{}
	rows, err := model.db.Query(
		"SELECT submission_id, form_id, created_at FROM submissions WHERE form_id = ? ORDER BY created_at ASC",
		formID,
	)
	if err != nil {
		return submissions, err
	}
	defer rows.Close()
	for rows.Next() {
		submission := Submission{}
		err := rows.Scan(&submission.id, &submission.formID, &submission.createAt)
		if err != nil {
			return submissions, err
		}
		submissions = append(submissions, submission)
	}
	err = rows.Err()
	return submissions, err
}

type sqlEntryModel struct {
	db *sql.DB
}

func (model sqlEntryModel) Create(submissionID, labelID int64, txt string) (Entry, error) {
	entry := Entry{labelID: labelID, submissionID: submissionID, txt: txt}
	stmt, err := model.db.Prepare("INSERT INTO entries (label_id, submission_id, txt) VALUES (?, ?, ?)")
	if err != nil {
		return Entry{}, err
	}
	res, err := stmt.Exec(entry.labelID, entry.submissionID, entry.txt)
	if err != nil {
		return Entry{}, err
	}
	entry.id, err = res.LastInsertId()
	if err != nil {
		return Entry{}, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func (model sqlEntryModel) GetEntries(submissionID, labelID int64) ([]Entry, error) {
	entries := []Entry{}
	rows, err := model.db.Query(
		`SELECT entry_id, submission_id, label_id, txt
		FROM entries WHERE submission_id = ? AND label_id = ?`,
		submissionID,
		labelID,
	)
	if err != nil {
		return entries, err
	}
	defer rows.Close()
	for rows.Next() {
		entry := Entry{}
		err := rows.Scan(&entry.id, &entry.submissionID, &entry.labelID, &entry.txt)
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	err = rows.Err()
	return entries, err
}
