// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type Entry struct {
	Date    string
	Content sql.NullString
	Keyword sql.NullString
	Mood    sql.NullString
	Remarks sql.NullString
}
