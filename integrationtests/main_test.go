package integrationtests

// Basic imports
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"alecgibson.ca/go-postgres-petstore/client"
	"github.com/stretchr/testify/suite"
)

const (
	defaultServerAddress = "http://localhost:5000"
	serverAddressKey     = "SERVER_ADDRESS"
	timeout              = 5 * time.Second
)

var (
	testPetName  = "fido"
	testPetName2 = "dover"
	testPetTag   = "dog"
	limitOne     = int32(1)
)

type IntegrationTestSuite struct {
	suite.Suite
	client client.ClientWithResponsesInterface
}

func (suite *IntegrationTestSuite) SetupTest() {
	serverAddress, found := os.LookupEnv(serverAddressKey)
	if !found {
		serverAddress = defaultServerAddress
	}

	fmt.Printf("Using address: %v\n", serverAddress)

	c, err := client.NewClientWithResponses(serverAddress)
	if err != nil {
		fmt.Println("Failed to create testing client")
		panic(err)
	}

	suite.client = c
}

func (suite *IntegrationTestSuite) TestCreateAndDeletePets() {
	addPetRequestBodies := []client.AddPetJSONRequestBody{{
		Name: testPetName,
		Tag:  &testPetTag,
	}, {
		Name: testPetName2,
		Tag:  nil,
	}}
	suite.performAddPetRequests(addPetRequestBodies)

	suite.performFindPetsAndAssertNumResults(client.FindPetsParams{}, len(addPetRequestBodies))
	suite.performFindPetsAndAssertNumResults(client.FindPetsParams{
		Tags: &[]string{testPetTag},
	}, 1)
	suite.performFindPetsAndAssertNumResults(client.FindPetsParams{
		Limit: &limitOne,
	}, 1)

	suite.checkFindPetByIDMatchesFindPets()
	suite.performDeletePetAndCheckPetWasDeleted()
}

// Performs a AddPet operation for each body provided
// Asserts that each pet returned in the responses match the fields in the request bodies
func (suite *IntegrationTestSuite) performAddPetRequests(requestBodies []client.AddPetJSONRequestBody) {
	for _, body := range requestBodies {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		resp, err := suite.client.AddPetWithResponse(ctx, body)
		suite.Assert().NoError(err)
		suite.Assert().NotNil(resp)
		suite.Assert().Equal(http.StatusOK, (*(*resp).HTTPResponse).StatusCode)
		suite.Assert().Equal(body.Name, (*(*resp).JSON200).Name)
		if body.Tag == nil {
			suite.Assert().Nil((*(*resp).JSON200).Tag)
		} else {
			suite.Assert().Equal(*body.Tag, *(*(*resp).JSON200).Tag)
		}
	}
}

// Performs a FindPets operation using the provided params
// Asserts that the number of pets in the response matches `numResults`
func (suite *IntegrationTestSuite) performFindPetsAndAssertNumResults(params client.FindPetsParams, numResults int) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := suite.client.FindPetsWithResponse(ctx, &params)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(http.StatusOK, (*(*resp).HTTPResponse).StatusCode)
	suite.Assert().Equal(numResults, len(*(*resp).JSON200))
}

// Performs a FindPets operation, and saves the first pet from the response
// Then performs a FindPetByID using the ID of the saved pet
// Asserts that all fields from the FindPetByID response match the saved pet
func (suite *IntegrationTestSuite) checkFindPetByIDMatchesFindPets() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	findPetsResp, err := suite.client.FindPetsWithResponse(ctx, &client.FindPetsParams{})
	suite.Assert().NoError(err)
	suite.Assert().NotNil(findPetsResp)
	suite.Assert().Equal(http.StatusOK, (*(*findPetsResp).HTTPResponse).StatusCode)
	suite.Assert().NotEqual(0, len(*(*findPetsResp).JSON200))

	pet := (*(*findPetsResp).JSON200)[0]

	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	findPetByIDResp, err := suite.client.FindPetByIDWithResponse(ctx, pet.Id)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(findPetByIDResp)
	suite.Assert().Equal(http.StatusOK, (*(*findPetByIDResp).HTTPResponse).StatusCode)

	petByID := *(*findPetByIDResp).JSON200
	suite.Assert().Equal(pet.Id, petByID.Id)
	suite.Assert().Equal(pet.Name, petByID.Name)
	suite.Assert().Equal(pet.Tag, petByID.Tag)
}

func (suite *IntegrationTestSuite) performDeletePetAndCheckPetWasDeleted() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	findPetsResp, err := suite.client.FindPetsWithResponse(ctx, &client.FindPetsParams{})
	suite.Assert().NoError(err)
	suite.Assert().NotNil(findPetsResp)
	suite.Assert().Equal(http.StatusOK, (*(*findPetsResp).HTTPResponse).StatusCode)
	suite.Assert().NotEqual(0, len(*(*findPetsResp).JSON200))

	pet := (*(*findPetsResp).JSON200)[0]

	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	deletePetResp, err := suite.client.DeletePetWithResponse(ctx, pet.Id)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(deletePetResp)
	suite.Assert().Equal(http.StatusNoContent, (*(*deletePetResp).HTTPResponse).StatusCode)

	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	findPetByIDResp, err := suite.client.FindPetByIDWithResponse(ctx, pet.Id)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(findPetByIDResp)
	suite.Assert().Equal(http.StatusNotFound, (*(*findPetByIDResp).HTTPResponse).StatusCode)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
