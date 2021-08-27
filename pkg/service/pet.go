package service

import (
	"context"
	"database/sql"

	"alecgibson.ca/go-postgres-petstore/pkg/db"
	"github.com/jackc/pgx/v4"
	"github.com/palantir/stacktrace"
)

type (
	Pet interface {
		FindAll(ctx context.Context, tags *[]string, limit *int32) ([]db.PetstorePet, error)
		Create(ctx context.Context, name string, tag *string) (db.PetstorePet, error)
		Delete(ctx context.Context, id int64) error
		FindByID(ctx context.Context, id int64) (db.PetstorePet, error)
	}

	pet struct {
		db db.Querier
	}
)

func (p *pet) FindAll(ctx context.Context, tagsPtr *[]string, limit *int32) ([]db.PetstorePet, error) {
	var pets []db.PetstorePet
	var err error

	tags := []string{}
	if tagsPtr != nil {
		tags = *tagsPtr
	}

	if limit != nil {
		params := db.ListPetsWithLimitParams{
			Tags:       tags,
			MaxRecords: *limit,
		}
		pets, err = p.db.ListPetsWithLimit(ctx, params)
		if err != nil {
			return nil, stacktrace.Propagate(err, "")
		}
	} else {
		pets, err = p.db.ListPets(ctx, tags)
		if err != nil {
			return nil, stacktrace.Propagate(err, "")
		}
	}

	return pets, nil
}

func (p *pet) Create(ctx context.Context, name string, tag *string) (db.PetstorePet, error) {
	tagParam := sql.NullString{}
	if tag != nil {
		tagParam.Valid = true
		tagParam.String = *tag
	}
	params := db.CreatePetParams{
		Name: name,
		Tag:  tagParam,
	}

	pet, err := p.db.CreatePet(ctx, params)
	if err != nil {
		return pet, stacktrace.Propagate(err, "")
	}

	return pet, nil
}

func (p *pet) Delete(ctx context.Context, id int64) error {
	err := p.db.DeletePet(ctx, id)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return stacktrace.PropagateWithCode(err, ErrNotFound, "")
		} else {
			return stacktrace.Propagate(err, "")
		}
	}

	return nil
}

func (p *pet) FindByID(ctx context.Context, id int64) (db.PetstorePet, error) {
	pet, err := p.db.FindPetByID(ctx, id)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return pet, stacktrace.PropagateWithCode(err, ErrNotFound, "")
		} else {
			return pet, stacktrace.Propagate(err, "")
		}
	}

	return pet, nil
}

func NewPetService(db db.Querier) Pet {
	return &pet{
		db: db,
	}
}
