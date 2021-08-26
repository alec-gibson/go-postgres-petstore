// Code generated by sqlc. DO NOT EDIT.
// source: queries.sql

package db

import (
	"context"
	"database/sql"
)

const CreatePet = `-- name: CreatePet :one
INSERT INTO pets (
  name, tag
) VALUES (
  $1, $2
)
RETURNING id, name, tag
`

type CreatePetParams struct {
	Name string         `db:"name"`
	Tag  sql.NullString `db:"tag"`
}

func (q *Queries) CreatePet(ctx context.Context, arg CreatePetParams) (Pet, error) {
	row := q.db.QueryRow(ctx, CreatePet, arg.Name, arg.Tag)
	var i Pet
	err := row.Scan(&i.ID, &i.Name, &i.Tag)
	return i, err
}

const DeletePet = `-- name: DeletePet :exec
DELETE FROM pets
WHERE id = $1
`

func (q *Queries) DeletePet(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, DeletePet, id)
	return err
}

const FindPetByID = `-- name: FindPetByID :one
SELECT id, name, tag FROM pets
WHERE id = $1 LIMIT 1
`

func (q *Queries) FindPetByID(ctx context.Context, id int64) (Pet, error) {
	row := q.db.QueryRow(ctx, FindPetByID, id)
	var i Pet
	err := row.Scan(&i.ID, &i.Name, &i.Tag)
	return i, err
}

const ListPets = `-- name: ListPets :many
SELECT id, name, tag FROM pets
WHERE cardinality($1::varchar[]) = 0 OR tag = ANY($1::varchar[])
ORDER BY name
`

func (q *Queries) ListPets(ctx context.Context, dollar_1 []string) ([]Pet, error) {
	rows, err := q.db.Query(ctx, ListPets, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pet
	for rows.Next() {
		var i Pet
		if err := rows.Scan(&i.ID, &i.Name, &i.Tag); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const ListPetsWithLimit = `-- name: ListPetsWithLimit :many
SELECT id, name, tag FROM pets
WHERE cardinality($1::varchar[]) = 0 OR tag = ANY($1::varchar[])
ORDER BY name LIMIT $2
`

type ListPetsWithLimitParams struct {
	Column1 []string `db:"column_1"`
	Limit   int32    `db:"limit"`
}

func (q *Queries) ListPetsWithLimit(ctx context.Context, arg ListPetsWithLimitParams) ([]Pet, error) {
	rows, err := q.db.Query(ctx, ListPetsWithLimit, arg.Column1, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pet
	for rows.Next() {
		var i Pet
		if err := rows.Scan(&i.ID, &i.Name, &i.Tag); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}