package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"alecgibson.ca/go-postgres-petstore/pkg/infrastructure/db"
	"alecgibson.ca/go-postgres-petstore/pkg/service"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/palantir/stacktrace"
	"github.com/stretchr/testify/suite"
)

var (
	testLimit       int32 = 5
	testTag               = "test"
	testID          int64 = 1
	testName              = "name"
	testErr               = stacktrace.NewError("testErr")
	testErrNotFound       = stacktrace.NewErrorWithCode(service.ErrNotFound, "testErr")
)

type ControllerImplSuite struct {
	suite.Suite
	ctx            echo.Context
	rec            *httptest.ResponseRecorder
	controller     ServerInterface
	mockPetService *service.MockPet
	allControllers []*gomock.Controller
}

func (suite *ControllerImplSuite) SetupTest() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	suite.rec = httptest.NewRecorder()
	suite.ctx = echo.New().NewContext(req, suite.rec)

	mockPetController := gomock.NewController(suite.T())
	suite.allControllers = append(suite.allControllers, mockPetController)

	suite.mockPetService = service.NewMockPet(mockPetController)
	suite.controller = NewController(suite.mockPetService)
}

func (suite *ControllerImplSuite) TearDownTest() {
	for _, controller := range suite.allControllers {
		controller.Finish()
	}
}

// Given: pet service returns a valid result for FindAll operation
// Test: controller returns a 200 status code and the correct response body
func (suite *ControllerImplSuite) TestFindPets() {
	findPetsResponse := []db.PetstorePet{{
		ID:   testID,
		Name: testName,
		Tag: sql.NullString{
			Valid:  true,
			String: testTag,
		},
	}}
	suite.mockPetService.EXPECT().FindAll(gomock.Any(), nil, nil).Return(findPetsResponse, nil)

	params := FindPetsParams{
		Tags:  nil,
		Limit: nil,
	}
	err := suite.controller.FindPets(suite.ctx, params)
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusOK, suite.rec.Result().StatusCode)

	responsePets := []Pet{}
	json.NewDecoder(suite.rec.Result().Body).Decode(&responsePets)
	suite.Assert().Equal(1, len(responsePets))
	suite.Assert().Equal(testID, responsePets[0].Id)
	suite.Assert().Equal(testName, responsePets[0].Name)
	suite.Assert().Equal(testTag, *responsePets[0].Tag)
}

// Given: pet service returns an error for FindAll operation
// Test: controller returns the same error
func (suite *ControllerImplSuite) TestFindPetsServerError() {
	suite.mockPetService.EXPECT().FindAll(gomock.Any(), nil, nil).Return(nil, testErr)

	params := FindPetsParams{
		Tags:  nil,
		Limit: nil,
	}
	err := suite.controller.FindPets(suite.ctx, params)
	suite.Assert().Error(err)
	suite.Assert().Equal(testErr.Error(), err.Error())
}

// Given: pet service returns a valid result for AddPet operation
// Test: controller returns a 200 status code and the correct response body
func (suite *ControllerImplSuite) TestAddPet() {
	addPetResponse := db.PetstorePet{
		ID:   testID,
		Name: testName,
		Tag: sql.NullString{
			Valid:  true,
			String: testTag,
		},
	}
	suite.mockPetService.EXPECT().Create(gomock.Any(), testName, &testTag).Return(addPetResponse, nil)

	addPetRequest := NewPet{
		Name: testName,
		Tag:  &testTag,
	}
	requestBody, _ := json.Marshal(addPetRequest)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.ctx.SetRequest(req)

	err := suite.controller.AddPet(suite.ctx)
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusOK, suite.rec.Result().StatusCode)

	responsePet := Pet{}
	json.NewDecoder(suite.rec.Result().Body).Decode(&responsePet)
	suite.Assert().Equal(testID, responsePet.Id)
	suite.Assert().Equal(testName, responsePet.Name)
	suite.Assert().Equal(testTag, *responsePet.Tag)
}

// Given: pet service returns an error for AddPet operation
// Test: controller returns the same error
func (suite *ControllerImplSuite) TestAddPetServerError() {
	suite.mockPetService.EXPECT().Create(gomock.Any(), testName, &testTag).Return(db.PetstorePet{}, testErr)

	addPetRequest := NewPet{
		Name: testName,
		Tag:  &testTag,
	}
	requestBody, _ := json.Marshal(addPetRequest)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.ctx.SetRequest(req)

	err := suite.controller.AddPet(suite.ctx)
	suite.Assert().Error(err)
	suite.Assert().Equal(testErr.Error(), err.Error())
}

// Given: pet service returns no error for Delete operation
// Test: controller returns a 204 status code
func (suite *ControllerImplSuite) TestDeletePet() {
	suite.mockPetService.EXPECT().Delete(gomock.Any(), testID).Return(nil)
	err := suite.controller.DeletePet(suite.ctx, testID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusNoContent, suite.rec.Result().StatusCode)
}

// Given: pet service returns an error for Delete operation
// Test: controller returns the same error
func (suite *ControllerImplSuite) TestDeletePetServerError() {
	suite.mockPetService.EXPECT().Delete(gomock.Any(), testID).Return(testErr)
	err := suite.controller.DeletePet(suite.ctx, testID)
	suite.Assert().Error(err)
	suite.Assert().Equal(testErr.Error(), err.Error())
}

// Given: pet service returns a valid response for FindByID operation
// Test: controller returns a 200 status code and the correct response body
func (suite *ControllerImplSuite) TestFindPetByID() {
	findPetResponse := db.PetstorePet{
		ID:   testID,
		Name: testName,
		Tag: sql.NullString{
			Valid:  true,
			String: testTag,
		},
	}
	suite.mockPetService.EXPECT().FindByID(gomock.Any(), testID).Return(findPetResponse, nil)
	err := suite.controller.FindPetByID(suite.ctx, testID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusOK, suite.rec.Result().StatusCode)

	responsePet := Pet{}
	json.NewDecoder(suite.rec.Result().Body).Decode(&responsePet)
	suite.Assert().Equal(testID, responsePet.Id)
	suite.Assert().Equal(testName, responsePet.Name)
	suite.Assert().Equal(testTag, *responsePet.Tag)
}

// Given: pet service returns ErrNotFound for FindByID operation
// Test: controller returns a 404 status code
func (suite *ControllerImplSuite) TestFindPetByIDNotFound() {
	suite.mockPetService.EXPECT().FindByID(gomock.Any(), testID).Return(db.PetstorePet{}, testErrNotFound)
	err := suite.controller.FindPetByID(suite.ctx, testID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusNotFound, suite.rec.Result().StatusCode)
}

func TestControllerImplSuite(t *testing.T) {
	suite.Run(t, new(ControllerImplSuite))
}
