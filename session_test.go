package main

import (
	"database/sql"
	"testing"
)

func TestGetEntriesByID(t *testing.T) {
	cases := []struct {
		desc   string
		task   session
		entrys []entry
		err    error
	}{
		{
			"valid ID which contains 1 valid session",
			session{id: 1, ksatID: 3},
			[]entry{
				{
					id:        1,
					sessionID: 1,
					promptID:  1,
					txt:       "first entry",
				},
				{
					id:        2,
					sessionID: 1,
					promptID:  2,
					txt:       "second entry",
				},
			},
			nil,
		},
		{
			"valid ID but it has no session",
			session{id: 2, ksatID: 3},
			[]entry{},
			nil,
		},
		{
			"valid id for a session that does not exist",
			session{id: 101010, ksatID: 3},
			nil,
			sql.ErrNoRows,
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			entrys, err := tc.task.getEntriesByID()
			if err != tc.err {
				t.Fatalf("errors don't match: %v, %v", err, tc.err)
			}
			for i, item := range entrys {
				// needed since I am creating more entries than are shown
				if i >= len(tc.entrys) {
					break
				}
				if item.id != tc.entrys[i].id {
					t.Fatalf("ids don't match: %v, %v", item.id, tc.entrys[i].id)
				}
				if item.sessionID != tc.entrys[i].sessionID {
					t.Fatalf("sessionIDs don't match: %v, %v", item.sessionID, tc.entrys[i].sessionID)
				}
				if item.promptID != tc.entrys[i].promptID {
					t.Fatalf("promptIDs don't match: %v, %v", item.promptID, tc.entrys[i].promptID)
				}
				if item.txt != tc.entrys[i].txt {
					t.Fatalf("txts don't match: %v, %v", item.txt, tc.entrys[i].txt)
				}
			}
		})
	}
}
