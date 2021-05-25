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
	ID          int64
	Name, Usage string
}

// FormModel ...
type FormModel interface {
	Create(name, usage string) (Form, error)
	GetByName(name string) (Form, error)
	GetByID(id int64) (Form, error)
	GetAll() ([]Form, error)
	DeleteByID(id int64) (Form, error)
	DeleteByName(name string) (Form, error)
	Update(formID int64, name, usage string) (Form, error)
}

// Label ...
type Label struct {
	ID, FormID, Position int64
	Repeatable           bool
	Name, Usage          string
}

// LabelModel ...
type LabelModel interface {
	Create(formID, position int64, repeatable bool, name, usage string) (Label, error)
	GetLabels(formID int64) ([]Label, error)
	Update(labelID int64, name, usage string) (Label, error)
}

// Submission ...
type Submission struct {
	ID, FormID int64
	CreateAt   time.Time
}

// SubmissionModel ...
type SubmissionModel interface {
	Create(formID int64) (Submission, error)
	GetSubmissions(formID int64) ([]Submission, error)
}

// Entry ...
type Entry struct {
	ID, LabelID, SubmissionID int64
	Txt                       string
}

// EntryModel ...
type EntryModel interface {
	Create(submissionID, labelID int64, txt string) (Entry, error)
	GetEntries(submissionID, labeID int64) ([]Entry, error)
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
