package formly

import (
	"errors"
	"regexp"
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
	Create(name, usage string) (Form, error)
	GetByName(name string) (Form, error)
	GetByID(id int64) (Form, error)
	GetAll() ([]Form, error)
	DeleteByName(name string) error
}

// ErrInvalidLengthName ...
var ErrInvalidLengthName error = errors.New("name's length is not between 1 - 16 characters long")

// ErrNameIsNotAWord ...
var ErrNameIsNotAWord error = errors.New("name is not composed only of letters")

// ValidateName ...
func ValidateName(name string) error {
	if !(len(name) >= 1 && len(name) <= 16) {
		return ErrInvalidLengthName
	}
	if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(name) {
		return ErrNameIsNotAWord
	}
	return nil
}

// ErrInvalidLengthUsage ...
var ErrInvalidLengthUsage error = errors.New("name's length is not between 5 - 252 characters long")

// ValidateUsage ...
func ValidateUsage(usage string) error {
	if !(len(usage) >= 5 && len(usage) <= 252) {
		return ErrInvalidLengthUsage
	}
	return nil
}

// Label ...
type Label struct {
	id, formID, position int64
	repeatable           bool
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

// GetRepeatable gets repeatable from Label
func (label Label) GetRepeatable() bool {
	return label.repeatable
}

// LabelModel ...
type LabelModel interface {
	Create(formID, position int64, repeatable bool, name, usage string) (Label, error)
	GetLabels(formID int64) ([]Label, error)
}

// Submission ...
type Submission struct {
	id, formID int64
	createAt   time.Time
}

// GetID gets id from Submission
func (submission Submission) GetID() int64 {
	return submission.id
}

// GetCreateAt gets createAt from Submission
func (submission Submission) GetCreateAt() time.Time {
	return submission.createAt
}

// SubmissionModel ...
type SubmissionModel interface {
	Create(formID int64) (Submission, error)
	GetSubmissions(formID int64) ([]Submission, error)
}

// Entry ...
type Entry struct {
	id, labelID, submissionID int64
	txt                       string
}

// GetID gets id from Entry
func (entry Entry) GetID() int64 {
	return entry.id
}

// GetTxt gets txt from Entry
func (entry Entry) GetTxt() string {
	return entry.txt
}

// EntryModel ...
type EntryModel interface {
	Create(submissionID, labelID int64, txt string) (Entry, error)
	GetEntries(submissionID, labeID int64) ([]Entry, error)
}
