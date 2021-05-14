package ksat

import (
	"database/sql"
)

type sqlDriver struct {
	db   *sql.DB
	form sqlFormDriver
}

func newSQLDriver(db *sql.DB) sqlDriver {
	return sqlDriver{db: db, form: sqlFormDriver{db: db}}
}

type sqlFormDriver struct {
	db *sql.DB
}

func (driver sqlFormDriver) GetByName(name string) (Form, error) {
	form := Form{}
	err := driver.db.QueryRow(
		"SELECT form_id, name, usage FROM forms WHERE name = ?",
		name,
	).Scan(&form.id, &form.name, &form.usage)
	if err == sql.ErrNoRows {
		err = nil
	}
	return form, err
}
func (driver sqlFormDriver) GetByID(id int64) (Form, error) {
	form := Form{}
	err := driver.db.QueryRow(
		"SELECT form_id, name, usage FROM forms WHERE id = ?",
		id,
	).Scan(&form.id, &form.name, &form.usage)
	if err == sql.ErrNoRows {
		err = nil
	}
	return form, err
}
func (driver sqlFormDriver) GetForms() ([]Form, error) {
	forms := []Form{}
	rows, err := driver.db.Query("SELECT form_id, name, usage FROM forms")
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
func (driver sqlFormDriver) GetLabels(formID int64) ([]Label, error) {
	labels := []Label{}
	rows, err := driver.db.Query("SELECT label_id, position, name, usage FROM labels where form_id = ? ORDER BY position ASC", formID)
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
