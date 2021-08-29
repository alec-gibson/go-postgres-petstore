package service

import (
	context "context"
	"database/sql"
	"testing"

	db "alecgibson.ca/go-postgres-petstore/pkg/infrastructure/db"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/palantir/stacktrace"
	"github.com/stretchr/testify/suite"
)

var (
	testLimit int32 = 5
	testTag         = "test"
	testID    int64 = 1
	testName        = "name"
)

type PetServiceSuite struct {
	suite.Suite
	ctx            context.Context
	petService     Pet
	mockQuerier    *db.MockQuerier
	allControllers []*gomock.Controller
}

func (suite *PetServiceSuite) SetupTest() {
	suite.ctx = context.Background()
	mockQuerierController := gomock.NewController(suite.T())
	suite.allControllers = append(suite.allControllers, mockQuerierController)

	suite.mockQuerier = db.NewMockQuerier(mockQuerierController)
	suite.petService = NewPetService(suite.mockQuerier)
}

func (suite *PetServiceSuite) TearDownTest() {
	for _, controller := range suite.allControllers {
		controller.Finish()
	}
}

func (suite *PetServiceSuite) TestFindAll() {
	testPets := []db.PetstorePet{{
		ID:   testID,
		Name: testName,
		Tag: sql.NullString{
			Valid:  true,
			String: testTag,
		},
	}}
	suite.mockQuerier.EXPECT().ListPetsWithLimit(suite.ctx, gomock.Any()).Return(testPets, nil)
	pets, err := suite.petService.FindAll(suite.ctx, nil, &testLimit)
	suite.Assert().NoError(err)
	suite.Assert().Equal(1, len(pets))
	suite.Assert().Equal(testID, pets[0].ID)
}

func (suite *PetServiceSuite) TestCreate() {
	testPet := db.PetstorePet{
		ID:   testID,
		Name: testName,
		Tag: sql.NullString{
			Valid:  true,
			String: testTag,
		},
	}
	suite.mockQuerier.EXPECT().CreatePet(suite.ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, params db.CreatePetParams) (db.PetstorePet, error) {
		suite.Assert().Equal(testName, params.Name)
		suite.Assert().Equal(testTag, params.Tag.String)
		return testPet, nil
	})

	pet, err := suite.petService.Create(suite.ctx, testName, &testTag)
	suite.Assert().NoError(err)
	suite.Assert().Equal(testPet.ID, pet.ID)
}

func (suite *PetServiceSuite) TestDelete() {
	suite.mockQuerier.EXPECT().DeletePet(suite.ctx, testID).Return(nil)
	err := suite.petService.Delete(suite.ctx, testID)
	suite.Assert().NoError(err)
}

func (suite *PetServiceSuite) TestDeleteNotFound() {
	suite.mockQuerier.EXPECT().DeletePet(suite.ctx, testID).Return(pgx.ErrNoRows)
	err := suite.petService.Delete(suite.ctx, testID)
	suite.Assert().Error(err)
	suite.Assert().Equal(ErrNotFound, stacktrace.GetCode(err))
}

func (suite *PetServiceSuite) TestFindByID() {
	testPet := db.PetstorePet{
		ID:   testID,
		Name: testName,
		Tag: sql.NullString{
			Valid:  true,
			String: testTag,
		},
	}
	suite.mockQuerier.EXPECT().FindPetByID(suite.ctx, testID).Return(testPet, nil)
	pet, err := suite.petService.FindByID(suite.ctx, testID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(testPet.ID, pet.ID)
}

func (suite *PetServiceSuite) TestFindByIDNotFound() {
	suite.mockQuerier.EXPECT().FindPetByID(suite.ctx, testID).Return(db.PetstorePet{}, pgx.ErrNoRows)
	_, err := suite.petService.FindByID(suite.ctx, testID)
	suite.Assert().Error(err)
	suite.Assert().Equal(ErrNotFound, stacktrace.GetCode(err))
}

func TestPetServiceSuite(t *testing.T) {
	suite.Run(t, new(PetServiceSuite))
}
