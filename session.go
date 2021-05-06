package main

import "time"

type session struct {
	id, ksatID int64
	createAt   time.Time
}

func (session *session) dbInsert() error {
	task := ksat{id: session.ksatID}
	if err := task.getByID(); err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO sessions (ksat_id) VALUES (?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(session.ksatID)
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
	).Scan(&session.id, &session.ksatID, &session.createAt)
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
