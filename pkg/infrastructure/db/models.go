// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
)

type PetstorePet struct {
	ID   int64          `db:"id"`
	Name string         `db:"name"`
	Tag  sql.NullString `db:"tag"`
}