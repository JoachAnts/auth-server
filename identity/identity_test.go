package identity_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JoachAnts/auth-server/identity"
	"github.com/JoachAnts/auth-server/repo"
	"github.com/stretchr/testify/assert"
)

type TestIdentityResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}

type Role struct {
	CompanyID string `json:"companyID"`
	Role      string `json:"role"`
}

func testHandler(t *testing.T, userID *string, expectedStatus int, expectedBody *TestIdentityResponse) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/me", nil)
	if userID != nil {
		req.Header.Add("authorization", *userID)
	}
	if err != nil {
		t.Fatal(err)
	}

	identity.NewHandler(repo.NewRepo(map[string]repo.User{
		"1": {
			ID:   "1",
			Name: "John Reece",
			Roles: []repo.Role{
				{
					CompanyID: "1",
					Role:      "user",
				},
				{
					CompanyID: "2",
					Role:      "admin",
				},
			},
		},
		"2": {
			ID:   "2",
			Name: "Bob Smith",
			Roles: []repo.Role{
				{
					CompanyID: "1",
					Role:      "admin",
				},
				{
					CompanyID: "2",
					Role:      "user",
				},
			},
		},
		"3": {
			ID:    "3",
			Name:  "Eve",
			Roles: nil,
		},
	}, map[string](map[string]repo.Card){})).ServeHTTP(writer, req)

	assert.Equal(t, expectedStatus, writer.Result().StatusCode)
	resBody := TestIdentityResponse{}
	if expectedBody == nil {
		return
	}
	if err = json.Unmarshal(writer.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}
	assert.Equal(t, expectedBody.ID, resBody.ID)
	assert.Equal(t, expectedBody.Name, resBody.Name)
	assert.Equal(t, expectedBody.Roles, resBody.Roles)
}

func TestIdentityUser1(t *testing.T) {
	userID := "1"
	testHandler(t, &userID, 200, &TestIdentityResponse{
		ID:   "1",
		Name: "John Reece",
		Roles: []Role{
			{
				CompanyID: "1",
				Role:      "user",
			},
			{
				CompanyID: "2",
				Role:      "admin",
			},
		},
	})
}

func TestIdentityUser2(t *testing.T) {
	userID := "2"
	testHandler(t, &userID, 200, &TestIdentityResponse{
		ID:   "2",
		Name: "Bob Smith",
		Roles: []Role{
			{
				CompanyID: "1",
				Role:      "admin",
			},
			{
				CompanyID: "2",
				Role:      "user",
			},
		},
	})
}

func TestIdentityUserWithoutRoles(t *testing.T) {
	userID := "3"
	testHandler(t, &userID, http.StatusForbidden, nil)
}

func TestIdentityUnauthenticated(t *testing.T) {
	testHandler(t, nil, http.StatusUnauthorized, nil)
}

func TestIdentityNotExisting(t *testing.T) {
	userID := "1001"
	testHandler(t, &userID, http.StatusNotFound, nil)
}
