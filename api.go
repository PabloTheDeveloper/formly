package main

import (
	"database/sql"
	"fmt"
)

type noKsatIdByNameErr struct {
	name string
	err  error
}

func (e noKsatIdByNameErr) Error() string {
	return e.err.Error()
}
func (e noKsatIdByNameErr) UserError() string {
	return fmt.Sprintf("'%v' does not exist", e.name)
}

// getKsatIdByName if returns 0 means item not found
func getKsatIdByName(name string) (id int64, err error) {
	err = isKsatNameValid(name)
	if err != nil {
		return
	}
	err = db.QueryRow("SELECT ksat_id FROM ksats WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		err = noKsatIdByNameErr{name: name, err: err}
		return
	}
	return
}
