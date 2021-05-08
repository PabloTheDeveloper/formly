package main

import "time"

type session struct {
	id, KsatID int64
	createAt   time.Time
}

func (session *session) dbInsert() error {
	if _, err := GetKsatByID(session.KsatID); err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO sessions (ksat_id) VALUES (?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(session.KsatID)
	if err != nil {
		return err
	}
	session.id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
func (session *session) getByID() error {
	return db.QueryRow(
		"SELECT session_id, ksat_id, created_at FROM sessions WHERE session_id = ?",
		session.id,
	).Scan(&session.id, &session.KsatID, &session.createAt)
}
func (session *session) getEntriesByID() ([]entry, error) {
	if err := session.getByID(); err != nil {
		return nil, err
	}
	entries := []entry{}
	rows, err := db.Query("SELECT entry_id, session_id, prompt_id, txt FROM entries WHERE session_id = ?", session.id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := entry{}
		err := rows.Scan(&item.id, &item.sessionID, &item.promptID, &item.txt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return entries, nil
}

type entry struct {
	id, sessionID, promptID int64
	txt                     string
}

func (entry *entry) dbInsert() error {
	session := session{id: entry.sessionID}
	if err := session.getByID(); err != nil {
		return err
	}
	prompt := prompt{id: entry.promptID}
	if err := prompt.getByID(); err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO entries (session_id, prompt_id, txt) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(entry.sessionID, entry.promptID, entry.txt)
	if err != nil {
		return err
	}
	entry.id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
