package ksat

import (
	"time"

	// to support sqlite
	_ "github.com/mattn/go-sqlite3"
)

// Env ...
type Env struct {
	FormModel
	LabelModel
	SubmissionModel
	EntryModel
	close func() error
}

// Close ...
func (env *Env) Close() error {
	return env.close()
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

// FormModel ...
type FormModel interface {
	GetByName(name string) (Form, error)
	GetByID(id int64) (Form, error)
	GetAll() ([]Form, error)
}

type Label struct {
	id, formID, position, repeatable int64
	name, usage                      string
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

// GetRepeatable gets repeatable from Label
func (label Label) GetRepeatable() bool {
	return label.repeatable != 0
}

// LabelModel ...
type LabelModel interface {
	Create(formID, position, repeatable int64, name, usage string) (Label, error)
	GetLabels(formID int64) ([]Label, error)
}

type Submission struct {
	id, formID int64
	createAt   time.Time
}

// GetID gets id from Submission
func (submission Submission) GetID() int64 {
	return submission.id
}

// SubmissionModel ...
type SubmissionModel interface {
	Create(formID int64) (Submission, error)
}

type Entry struct {
	id, labelID, submissionID int64
	txt                       string
}

// GetID gets id from Entry
func (entry Entry) GetID() int64 {
	return entry.id
}

// EntryModel ...
type EntryModel interface {
	Create(submissionID, labelID int64, txt string) (Entry, error)
}
