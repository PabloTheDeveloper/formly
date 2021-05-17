package ksat

import (
	"database/sql"
)

type sqlModels struct {
	db *sql.DB
	sqlFormModel
	sqlSubmissionModel
	sqlEntryModel
}

func newSQLEnv(db *sql.DB) *Env {
	return &Env{
		FormModel:       sqlFormModel{db: db},
		LabelModel:      sqlLabelModel{db: db},
		SubmissionModel: sqlSubmissionModel{db: db},
		EntryModel:      sqlEntryModel{db: db},
		close: func() error {
			return db.Close()
		},
	}
}

type sqlFormModel struct {
	db *sql.DB
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

func (model sqlLabelModel) Create(formID, position int64, name, usage string) (Label, error) {
	label := Label{formID: formID, position: position, name: name, usage: usage}
	stmt, err := model.db.Prepare(
		"INSERT INTO labels (form_id, position, name, usage) VALUES(?, ?, ?, ?)",
	)
	if err != nil {
		return Label{}, err
	}
	res, err := stmt.Exec(label.formID, label.position, label.name, label.usage)
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
	rows, err := model.db.Query("SELECT label_id, position, name, usage FROM labels where form_id = ? ORDER BY position ASC", formID)
	if err != nil {
		return labels, err
	}
	defer rows.Close()
	for rows.Next() {
		label := Label{}
		err := rows.Scan(&label.id, &label.position, &label.name, &label.usage)
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
