package card_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JoachAnts/auth-server/card"
	"github.com/JoachAnts/auth-server/repo"
	"github.com/stretchr/testify/assert"
)

type TestCardResponse struct {
	// TODO think about how to support different currencies
	MaskedNumber string `json:"maskedNumber"`
	Exp          string `json:"exp"`
	Limit        int    `json:"limit"`
	Balance      int    `json:"balance"`
}

func testGetHandler(t *testing.T, authUser *string, expectedStatus int, expectedBody *TestCardResponse) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/card", nil)
	if authUser != nil {
		req.Header.Add("authorization", *authUser)
	}
	if err != nil {
		t.Fatal(err)
	}

	card.NewHandler(repo.NewRepo(
		map[string]repo.User{},
		map[string]repo.Card{
			"1": {
				MaskedNumber: "**** **** **** 1111",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
		},
	)).ServeHTTP(writer, req)

	assert.Equal(t, expectedStatus, writer.Result().StatusCode)
	if expectedBody == nil {
		return
	}
	resBody := TestCardResponse{}
	if err = json.Unmarshal(writer.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}
	assert.Equal(t, expectedBody.MaskedNumber, resBody.MaskedNumber)
	assert.Equal(t, expectedBody.Exp, resBody.Exp)
	assert.Equal(t, expectedBody.Limit, resBody.Limit)
	assert.Equal(t, expectedBody.Balance, resBody.Balance)
}

func TestCardUser(t *testing.T) {
	userID := "1"
	testGetHandler(t, &userID, http.StatusOK, &TestCardResponse{
		MaskedNumber: "**** **** **** 1111",
		Exp:          "12/23",
		Limit:        10000,
		Balance:      4321,
	})
}

func TestCardUnauthorized(t *testing.T) {
	testGetHandler(t, nil, http.StatusUnauthorized, nil)
}

func TestNoCardForUser(t *testing.T) {
	userID := "1001"
	testGetHandler(t, &userID, http.StatusNotFound, nil)
}

type TestLimitRequest struct {
	UserID   string
	NewLimit int
}

func testChangeLimitAPI(t *testing.T, reqB *TestLimitRequest, authUser *string, expectedStatus int) {
	writer := httptest.NewRecorder()
	jBody, err := json.Marshal(reqB)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}
	req, err := http.NewRequest("POST", "/card", bytes.NewBuffer(jBody))
	if authUser != nil {
		req.Header.Add("authorization", *authUser)
	}
	if err != nil {
		t.Fatal(err)
	}

	card.NewHandler(repo.NewRepo(
		map[string]repo.User{
			"1": {
				ID:   "1",
				Name: "John Reece",
				Role: "user",
			},
			"2": {
				ID:   "2",
				Name: "Bob Smith",
				Role: "admin",
			},
		},
		map[string]repo.Card{
			"1": {
				MaskedNumber: "**** **** **** 1111",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
		},
	)).ServeHTTP(writer, req)

	assert.Equal(t, expectedStatus, writer.Result().StatusCode)
}

func TestUserShouldNotBeAbleToModifyLimit(t *testing.T) {
	userID := "1"
	testChangeLimitAPI(t, &TestLimitRequest{
		UserID:   "1",
		NewLimit: 20000,
	}, &userID, http.StatusForbidden)
}

func TestAdminShouldBeAbleToModifyLimit(t *testing.T) {
	userID := "2"
	testChangeLimitAPI(t, &TestLimitRequest{
		UserID:   "1",
		NewLimit: 20000,
	}, &userID, http.StatusOK)
}
