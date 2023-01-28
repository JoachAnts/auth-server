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

func TestCardUser(t *testing.T) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/card", nil)
	req.Header.Add("authorization", "1")
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

	assert.Equal(t, http.StatusOK, writer.Result().StatusCode)
	resBody := TestCardResponse{}
	if err = json.Unmarshal(writer.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}
	assert.Equal(t, "**** **** **** 1111", resBody.MaskedNumber)
	assert.Equal(t, "12/23", resBody.Exp)
	assert.Equal(t, 10000, resBody.Limit)
	assert.Equal(t, 4321, resBody.Balance)
}

func TestCardUnauthorized(t *testing.T) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/card", nil)
	if err != nil {
		t.Fatal(err)
	}

	card.NewHandler(repo.NewRepo(
		map[string]repo.User{},
		map[string]repo.Card{},
	)).ServeHTTP(writer, req)

	assert.Equal(t, http.StatusUnauthorized, writer.Result().StatusCode)
}

func TestNoCardForUser(t *testing.T) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/card", nil)
	req.Header.Add("authorization", "1")
	if err != nil {
		t.Fatal(err)
	}

	card.NewHandler(repo.NewRepo(
		map[string]repo.User{},
		map[string]repo.Card{},
	)).ServeHTTP(writer, req)

	assert.Equal(t, http.StatusNotFound, writer.Result().StatusCode)
}

type TestLimitRequest struct {
	UserID   string
	NewLimit int
}

func TestUserShouldNotBeAbleToModifyLimit(t *testing.T) {
	writer := httptest.NewRecorder()
	reqB := TestLimitRequest{
		UserID:   "1",
		NewLimit: 20000,
	}
	jBody, err := json.Marshal(reqB)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}
	req, err := http.NewRequest("POST", "/card", bytes.NewBuffer(jBody))
	req.Header.Add("authorization", "1")
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

	assert.Equal(t, http.StatusForbidden, writer.Result().StatusCode)
}

func TestAdminShouldBeAbleToModifyLimit(t *testing.T) {
	writer := httptest.NewRecorder()
	reqB := TestLimitRequest{
		UserID:   "1",
		NewLimit: 20000,
	}
	jBody, err := json.Marshal(reqB)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}
	req, err := http.NewRequest("POST", "/card", bytes.NewBuffer(jBody))
	req.Header.Add("authorization", "2")
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

	assert.Equal(t, http.StatusOK, writer.Result().StatusCode)
	resBody := TestCardResponse{}
	if err = json.Unmarshal(writer.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}
	assert.Equal(t, "**** **** **** 1111", resBody.MaskedNumber)
	assert.Equal(t, "12/23", resBody.Exp)
	assert.Equal(t, 20000, resBody.Limit)
	assert.Equal(t, 4321, resBody.Balance)
}
