package main

type session struct {
	id, ksatID int64
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
