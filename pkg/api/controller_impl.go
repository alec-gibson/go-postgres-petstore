package api

import (
	"context"
	"fmt"
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

func (c *controller) GetHealth(e echo.Context) error {
	return nil
}

func (c *controller) FindPets(e echo.Context, params FindPetsParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pets, err := c.petService.FindAll(ctx, params.Tags, params.Limit)
	if err != nil {
		fmt.Println(err)
		return respondByErrorCode(e, err)
	}

	response := []Pet{}
	for _, pet := range pets {
		responsePet := Pet{
			NewPet: NewPet{
				Name: pet.Name,
			},
			Id: pet.ID,
		}
		if pet.Tag.Valid {
			responsePet.Tag = &pet.Tag.String
		} else {
			responsePet.Tag = nil
		}

		response = append(response, responsePet)
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
		fmt.Println(err)
		return respondByErrorCode(e, err)
	}

	response := Pet{
		NewPet: NewPet{
			Name: createdPet.Name,
		},
		Id: createdPet.ID,
	}
	if createdPet.Tag.Valid {
		response.Tag = &createdPet.Tag.String
	} else {
		response.Tag = nil
	}

	e.JSON(http.StatusOK, response)
	return nil
}

func (c *controller) DeletePet(e echo.Context, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := c.petService.Delete(ctx, id)
	if err != nil {
		fmt.Println(err)
		return respondByErrorCode(e, err)
	}

	e.String(http.StatusNoContent, "")
	return nil
}

func (c *controller) FindPetByID(e echo.Context, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pet, err := c.petService.FindByID(ctx, id)
	if err != nil {
		fmt.Println(err)
		return respondByErrorCode(e, err)
	}

	response := Pet{
		NewPet: NewPet{
			Name: pet.Name,
		},
		Id: pet.ID,
	}
	if pet.Tag.Valid {
		response.Tag = &pet.Tag.String
	} else {
		response.Tag = nil
	}

	e.JSON(http.StatusOK, response)
	return nil
}

func respondByErrorCode(e echo.Context, err error) error {
	code := stacktrace.GetCode(err)
	switch code {
	case service.ErrNotFound:
		e.String(http.StatusNotFound, "Not found.")
	default:
		return err
	}
	return nil
}

func NewController(petService service.Pet) ServerInterface {
	return &controller{
		petService: petService,
	}
}
