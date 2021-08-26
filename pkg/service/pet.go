package service

import (
	"context"
	"database/sql"

	"alecgibson.ca/go-postgres-petstore/pkg/db"
	"github.com/palantir/stacktrace"
)

type (
	Pet interface {
		FindAll(ctx context.Context, tags *[]string, limit *int32) ([]db.Pet, error)
		Create(ctx context.Context, name string, tag *string) (db.Pet, error)
		Delete(ctx context.Context, id int64) error
		FindByID(ctx context.Context, id int64) (db.Pet, error)
	}

	pet struct {
		db db.Querier
	}
)

func (p *pet) FindAll(ctx context.Context, tagsPtr *[]string, limit *int32) ([]db.Pet, error) {
	var pets []db.Pet
	var err error

	tags := []string{}
	if tagsPtr != nil {
		tags = *tagsPtr
	}

	if limit != nil {
		params := db.ListPetsWithLimitParams{
			Column1: tags,
			Limit:   *limit,
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

func (p *pet) Create(ctx context.Context, name string, tag *string) (db.Pet, error) {
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
		return stacktrace.Propagate(err, "")
	}

	return nil
}

func (p *pet) FindByID(ctx context.Context, id int64) (db.Pet, error) {
	pet, err := p.db.FindPetByID(ctx, id)
	if err != nil {
		return pet, stacktrace.Propagate(err, "")
	}

	return pet, nil
}

func NewPetService(db db.Querier) Pet {
	return &pet{
		db: db,
	}
}