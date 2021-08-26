package api

import (
	"context"
	"net/http"
	"time"

	"alecgibson.ca/go-postgres-petstore/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/palantir/stacktrace"
)

const (
	timeout = 5 * time.Second
)

type controller struct {
	petService service.Pet
}

func (c *controller) FindPets(e echo.Context, params FindPetsParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pets, err := c.petService.FindAll(ctx, params.Tags, params.Limit)
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	response := []Pet{}
	for _, pet := range pets {
		response = append(response, Pet{
			NewPet: NewPet{
				Name: pet.Name,
				Tag:  &pet.Tag.String,
			},
			Id: pet.ID,
		})
	}
	e.JSON(http.StatusOK, response)
	return nil
}

func (c *controller) AddPet(e echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var pet NewPet
	if err := (&echo.DefaultBinder{}).BindBody(e, &pet); err != nil {
		return stacktrace.Propagate(err, "")
	}

	createdPet, err := c.petService.Create(ctx, pet.Name, pet.Tag)
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	response := Pet{
		NewPet: NewPet{
			Name: createdPet.Name,
			Tag:  &createdPet.Tag.String,
		},
		Id: createdPet.ID,
	}
	e.JSON(http.StatusOK, response)
	return nil
}

func (c *controller) DeletePet(e echo.Context, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := c.petService.Delete(ctx, id)
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	e.Response().Status = http.StatusNoContent
	return nil
}

func (c *controller) FindPetByID(e echo.Context, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pet, err := c.petService.FindByID(ctx, id)
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	response := Pet{
		NewPet: NewPet{
			Name: pet.Name,
			Tag:  &pet.Tag.String,
		},
		Id: pet.ID,
	}
	e.JSON(http.StatusOK, response)
	return nil
}

func NewController(petService service.Pet) ServerInterface {
	return &controller{
		petService: petService,
	}
}
