package formly

import (
	"database/sql"
	"fmt"
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
			name TEXT UNIQUE NOT NULL CHECK(length(name) >= 1 AND length(name) <= 16),
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
		`
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dataPath := filepath.Join(homePath, ".local", "share", "formly")
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
	form := Form{Name: name, Usage: usage}
	if err := model.db.QueryRow(
		"INSERT INTO forms (name, usage) VALUES (?, ?) RETURNING form_id",
		name,
		usage,
	).Scan(&form.ID); err != nil {
		return Form{}, err
	}
	return form, nil
}
func (model sqlFormModel) GetByName(name string) (Form, error) {
	form := Form{}
	if err := model.db.QueryRow(
		"SELECT form_id, name, usage FROM forms WHERE name = ?",
		name,
	).Scan(&form.ID, &form.Name, &form.Usage); err != nil {
		return Form{}, err
	}
	return form, nil
}
func (model sqlFormModel) GetByID(id int64) (Form, error) {
	form := Form{}
	if err := model.db.QueryRow(
		"SELECT form_id, name, usage FROM forms WHERE form_id = ?",
		id,
	).Scan(&form.ID, &form.Name, &form.Usage); err != nil {
		return Form{}, err
	}
	return form, nil
}
func (model sqlFormModel) GetAll() ([]Form, error) {
	forms := []Form{}
	rows, err := model.db.Query("SELECT form_id, name, usage FROM forms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		form := Form{}
		if err := rows.Scan(&form.ID, &form.Name, &form.Usage); err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return forms, nil
}
func (model sqlFormModel) DeleteByID(id int64) (Form, error) {
	form, err := model.GetByID(id)
	if err != nil {
		return Form{}, err
	}
	if _, err := model.db.Exec("DELETE FROM forms WHERE form_id = ?", id); err != nil {
		return Form{}, err
	}
	return form, nil
}

func (model sqlFormModel) DeleteByName(name string) (Form, error) {
	form, err := model.GetByName(name)
	if err != nil {
		return Form{}, err
	}
	if _, err := model.db.Exec("DELETE FROM forms WHERE name= ?", name); err != nil {
		return Form{}, err
	}
	return form, nil
}

func (model sqlFormModel) Update(formID int64, name, usage string) (Form, error) {
	if _, err := model.db.Exec(
		"UPDATE forms SET name = ?, usage = ? WHERE form_id = ?",
		name,
		usage,
		formID,
	); err != nil {
		return Form{}, err
	}
	return Form{ID: formID, Name: name, Usage: usage}, nil
}

type sqlLabelModel struct {
	db *sql.DB
}

func (model sqlLabelModel) Create(formID, position int64, repeatable bool, name, usage string) (Label, error) {
	formModel := sqlFormModel{db: model.db}
	if _, err := formModel.GetByID(formID); err != nil {
		return Label{}, err
	}
	label := Label{FormID: formID, Position: position, Name: name, Usage: usage}
	if err := model.db.QueryRow(
		"INSERT INTO labels (form_ID, position, repeatable, Name, Usage) VALUES (?, ?, ?, ?, ?) RETURNING label_id",
		formID,
		position,
		repeatable,
		name,
		usage,
	).Scan(&label.ID); err != nil {
		return Label{}, err
	}
	return Label{FormID: formID, Position: position, Repeatable: repeatable, Name: name, Usage: usage}, nil
}

func (model sqlLabelModel) GetLabels(formID int64) ([]Label, error) {
	formModel := sqlFormModel{db: model.db}
	if form, err := formModel.GetByID(formID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("form '%s' does not exists", form.Name)
		}
		return nil, err
	}
	labels := []Label{}
	rows, err := model.db.Query(
		"SELECT label_id, position, repeatable, name, usage FROM labels WHERE form_id = ? ORDER BY position ASC", formID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		label := Label{}
		if err := rows.Scan(&label.ID, &label.Position, &label.Repeatable, &label.Name, &label.Usage); err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return labels, nil
}

func (model sqlLabelModel) Update(formID, labelID, position int64, repeatable bool, name, usage string) ([]Label, error) {
	labels, err := model.GetLabels(formID)
	if err != nil {
		return nil, err
	}
	if int(position) > len(labels) || int(position) < 1 {
		return nil, fmt.Errorf("position has to be in range between: %v - %v", 1, len(labels))
	}
	// check that only name matches with the label with the same labelID
	var updatingLabel Label
	var swapLabel Label
	for _, label := range labels {
		if label.ID == labelID {
			updatingLabel = label
		}
		if label.ID != labelID && label.Name == name {
			return nil, fmt.Errorf("suggested new label name '%s' already exists", name)
		}
		if label.Position == position {
			swapLabel = label
		}
	}
	if _, err := model.db.Exec(
		"UPDATE labels SET name = ?, usage = ?, repeatable = ?, position = ? WHERE label_id = ? ",
		name,
		usage,
		repeatable,
		position,
		labelID,
	); err != nil {
		return nil, err
	}
	swapLabel.Position = updatingLabel.Position
	updatingLabel.Name = name
	updatingLabel.Usage = usage
	updatingLabel.Repeatable = repeatable
	if swapLabel.ID == updatingLabel.ID {
		updatingLabel.Position = position
		return []Label{updatingLabel}, nil
	}
	if _, err := model.db.Exec(
		"UPDATE labels SET position = ? WHERE label_id = ? ",
		swapLabel.Position,
		swapLabel.ID,
	); err != nil {
		return nil, err
	}

	updatingLabel.Position = position
	return []Label{updatingLabel, swapLabel}, nil
}

type sqlSubmissionModel struct {
	db *sql.DB
}

func (model sqlSubmissionModel) Create(formID int64) (Submission, error) {
	formModel := sqlFormModel{db: model.db}
	if form, err := formModel.GetByID(formID); err != nil {
		if err == sql.ErrNoRows {
			return Submission{}, fmt.Errorf("form '%s' does not exists", form.Name)
		}
		return Submission{}, err
	}
	submission := Submission{FormID: formID}
	if err := model.db.QueryRow(
		"INSERT INTO submissions (form_id) VALUES (?) RETURNING submission_id",
		formID,
	).Scan(&submission.ID); err != nil {
		return Submission{}, err
	}
	return submission, nil
}
func (model sqlSubmissionModel) GetSubmissions(formID int64) ([]Submission, error) {
	formModel := sqlFormModel{db: model.db}
	if form, err := formModel.GetByID(formID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("form '%s' does not exists", form.Name)
		}
		return nil, err
	}
	submissions := []Submission{}
	rows, err := model.db.Query(
		"SELECT submission_id, form_id, created_at FROM submissions WHERE form_id = ? ORDER BY created_at ASC",
		formID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		submission := Submission{}
		err := rows.Scan(&submission.ID, &submission.FormID, &submission.CreateAt)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return submissions, err
}

type sqlEntryModel struct {
	db *sql.DB
}

func (model sqlEntryModel) Create(submissionID, labelID int64, txt string) (Entry, error) {
	entry := Entry{LabelID: labelID, SubmissionID: submissionID, Txt: txt}
	if err := model.db.QueryRow(
		"INSERT INTO entries (label_id, submission_id, txt) VALUES (?, ?, ?) RETURNING entry_id",
		labelID,
		submissionID,
		txt,
	).Scan(&entry.ID); err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func (model sqlEntryModel) GetEntries(submissionID, labelID int64) ([]Entry, error) {
	entries := []Entry{}
	rows, err := model.db.Query(
		`SELECT entry_id, submission_id, label_id, txt FROM entries WHERE submission_id = ? AND label_id = ?`,
		submissionID,
		labelID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		entry := Entry{}
		err := rows.Scan(&entry.ID, &entry.SubmissionID, &entry.LabelID, &entry.Txt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	err = rows.Err()
	return entries, err
}
