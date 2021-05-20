package ksat

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func compareErr(returned, expected error) error {
	if expected != returned {
		return fmt.Errorf(
			`expected and returned err don't match.
			returned err: %v
			expected err: %v`,
			returned,
			expected,
		)
	}
	return nil
}

func compareStructures(returned, expected interface{}) error {
	if !reflect.DeepEqual(returned, expected) {
		return fmt.Errorf(
			`expected and returned %[1]s don't match.
			returned %[1]s: %[2]v
			expected %[1]s: %[3]v`,
			reflect.TypeOf(returned).Name(),
			returned,
			expected,
		)
	}
	return nil
}

func TestGetByName(t *testing.T) {
	type args struct {
		name string
	}
	type output struct {
		st  Form
		err error
	}
	type mockedCase struct {
		desc string
		args
		expected      output
		mockExpection func(args args, out output, mock sqlmock.Sqlmock)
	}
	sqlCmd := "SELECT form_id, name, usage FROM forms WHERE name = ?"
	columns := []string{"form_id", "name", "usage"}
	cases := []mockedCase{
		{
			desc:     "form with 'apple' name found and returned successfully",
			args:     args{"apple"},
			expected: output{Form{1, "apple", "some usage"}, nil},
			mockExpection: func(args args, out output, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).AddRow(out.st.id, out.st.name, out.st.usage)
				mock.ExpectQuery(sqlCmd).WithArgs(args.name).WillReturnRows(rows)
			},
		},
		{
			desc:     "form with 'ja a' name is not found. returns null case.",
			args:     args{"ja a"},
			expected: output{Form{}, nil},
			mockExpection: func(args args, out output, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(sqlCmd).WithArgs(args.name).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			desc:     `while looking for form with name 'ze', a hurricane strikes. return random err.`,
			args:     args{"ze"},
			expected: output{Form{}, errors.New("cable connecting db broke")},
			mockExpection: func(args args, out output, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(sqlCmd).WithArgs(args.name).WillReturnError(out.err)
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// setup sqlmock + env
			_db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal(err)
			}
			defer _db.Close()
			sqlDriver := newSQLDriver(_db)
			tc.mockExpection(tc.args, tc.expected, mock)

			ret := output{}
			ret.st, ret.err = sqlDriver.form.GetByName(tc.args.name)

			if err := compareErr(ret.err, tc.expected.err); err != nil {
				t.Fatal(err)
			}
			if err := compareStructures(ret.st, tc.expected.st); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestGetByID(t *testing.T) {
	type args struct {
		id int64
	}
	type output struct {
		st  Form
		err error
	}
	type mockedCase struct {
		desc string
		args
		expected      output
		mockExpection func(args args, out output, mock sqlmock.Sqlmock)
	}
	sqlCmd := "SELECT form_id, name, usage FROM forms WHERE id = ?"
	columns := []string{"form_id", "name", "usage"}
	cases := []mockedCase{
		{
			desc:     "form with id '1' found and returned successfully",
			args:     args{1},
			expected: output{Form{1, "apple", "some usage"}, nil},
			mockExpection: func(args args, out output, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).AddRow(out.st.id, out.st.name, out.st.usage)
				mock.ExpectQuery(sqlCmd).WithArgs(args.id).WillReturnRows(rows)
			},
		},
		{
			desc:     "form with id '2' is not found. returns null case.",
			args:     args{2},
			expected: output{Form{}, nil},
			mockExpection: func(args args, out output, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(sqlCmd).WithArgs(args.id).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			desc:     `while looking for form with id '0', a hurricane strikes. return random err.`,
			args:     args{0},
			expected: output{Form{}, errors.New("cable connecting db broke")},
			mockExpection: func(args args, out output, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(sqlCmd).WithArgs(args.id).WillReturnError(out.err)
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// setup sqlmock + env
			_db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal(err)
			}
			defer _db.Close()
			sqlDriver := newSQLDriver(_db)
			tc.mockExpection(tc.args, tc.expected, mock)

			ret := output{}
			ret.st, ret.err = sqlDriver.form.GetByID(tc.args.id)

			if err := compareErr(ret.err, tc.expected.err); err != nil {
				t.Fatal(err)
			}
			if err := compareStructures(ret.st, tc.expected.st); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestGetForms(t *testing.T) {
	type output struct {
		st  []Form
		err error
	}
	type mockedCase struct {
		desc          string
		expected      output
		mockExpection func(out output, mock sqlmock.Sqlmock)
	}
	sqlCmd := "SELECT form_id, name, usage FROM forms"
	columns := []string{"form_id", "name", "usage"}
	cases := []mockedCase{
		{
			desc: "2 forms returned successfully",
			expected: output{[]Form{
				{1, "apple", "some usage"},
				{2, "songs", "some usage"},
			}, nil},
			mockExpection: func(out output, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).
					AddRow(out.st[0].id, out.st[0].name, out.st[0].usage).
					AddRow(out.st[1].id, out.st[1].name, out.st[1].usage)
				mock.ExpectQuery(sqlCmd).WillReturnRows(rows)
			},
		},
		{
			desc:     "0 forms returned",
			expected: output{[]Form{}, nil},
			mockExpection: func(out output, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns)
				mock.ExpectQuery(sqlCmd).WillReturnRows(rows)
			},
		},
		{
			desc:     `a hurricane strikes. return random err.`,
			expected: output{[]Form{}, errors.New("cable connecting db broke")},
			mockExpection: func(out output, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(sqlCmd).WillReturnError(out.err)
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// setup sqlmock + env
			_db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal(err)
			}
			defer _db.Close()
			sqlDriver := newSQLDriver(_db)
			tc.mockExpection(tc.expected, mock)

			ret := output{}
			ret.st, ret.err = sqlDriver.form.GetForms()
			if err := compareErr(ret.err, tc.expected.err); err != nil {
				t.Fatal(err)
			}
			if err := compareStructures(ret.st, tc.expected.st); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestGetLabels(t *testing.T) {
	type args struct {
		formID int64
	}
	type output struct {
		st  []Label
		err error
	}
	type mockedCase struct {
		desc string
		args
		expected      output
		mockExpection func(args args, out output, mock sqlmock.Sqlmock)
	}
	sqlCmd := "SELECT label_id, position, name, usage FROM labels"
	columns := []string{"label_id", "position", "name", "usage"}
	cases := []mockedCase{
		{
			desc: "2 labels returned successfully",
			args: args{1},
			expected: output{[]Label{
				{1, 1, "apple", "some usage"},
				{2, 2, "songs", "some usage"},
			}, nil},
			mockExpection: func(args args, out output, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).
					AddRow(out.st[0].id, out.st[0].position, out.st[0].name, out.st[0].usage).
					AddRow(out.st[1].id, out.st[1].position, out.st[1].name, out.st[1].usage)
				mock.ExpectQuery(sqlCmd).WithArgs(args.formID).WillReturnRows(rows)
			},
		},
		{
			desc:     "0 labels returned",
			args:     args{1},
			expected: output{[]Label{}, nil},
			mockExpection: func(args args, out output, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns)
				mock.ExpectQuery(sqlCmd).WithArgs(args.formID).WillReturnRows(rows)
			},
		},
		{
			desc:     `a hurricane strikes. return random err.`,
			args:     args{1},
			expected: output{[]Label{}, errors.New("cable connecting db broke")},
			mockExpection: func(args args, out output, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(sqlCmd).WithArgs(args.formID).WillReturnError(out.err)
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// setup sqlmock + env
			_db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal(err)
			}
			defer _db.Close()
			sqlDriver := newSQLDriver(_db)
			tc.mockExpection(tc.args, tc.expected, mock)

			ret := output{}
			ret.st, ret.err = sqlDriver.form.GetLabels(tc.args.formID)
			if err := compareErr(ret.err, tc.expected.err); err != nil {
				t.Fatal(err)
			}
			if err := compareStructures(ret.st, tc.expected.st); err != nil {
				t.Fatal(err)
			}
		})
	}
}
