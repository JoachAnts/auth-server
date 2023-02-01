package card_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JoachAnts/auth-server/card"
	"github.com/JoachAnts/auth-server/repo"
	"github.com/stretchr/testify/assert"
)

type TestCardResponse struct {
	Cards []TestCard `json:"cards"`
}

type TestCard struct {
	CompanyID string          `json:"companyID"`
	Card      TestCardDetails `json:"card"`
}

type TestCardDetails struct {
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
		map[string](map[string]repo.Card){
			"1": {
				"1": {
					MaskedNumber: "**** **** **** 1111",
					Exp:          "12/23",
					Limit:        10000,
					Balance:      4321,
				},
				"2": {
					MaskedNumber: "**** **** **** 2222",
					Exp:          "12/23",
					Limit:        10000,
					Balance:      4321,
				},
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
	assert.Equal(t, *expectedBody, resBody)
}

func TestCardUser(t *testing.T) {
	userID := "1"
	testGetHandler(t, &userID, http.StatusOK, &TestCardResponse{
		Cards: []TestCard{
			{
				CompanyID: "1",
				Card: TestCardDetails{
					MaskedNumber: "**** **** **** 1111",
					Exp:          "12/23",
					Limit:        10000,
					Balance:      4321,
				},
			},
			{
				CompanyID: "2",
				Card: TestCardDetails{
					MaskedNumber: "**** **** **** 2222",
					Exp:          "12/23",
					Limit:        10000,
					Balance:      4321,
				},
			},
		},
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
	UserID    string
	CompanyID string
	NewLimit  int
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

	cards := map[string](map[string]repo.Card){
		"1": {
			"1": {
				MaskedNumber: "**** **** **** 1111",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
			"2": {
				MaskedNumber: "**** **** **** 2222",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
		},
		"2": {
			"1": {
				MaskedNumber: "**** **** **** 3333",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
			"2": {
				MaskedNumber: "**** **** **** 4444",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
		},
	}

	card.NewHandler(repo.NewRepo(
		map[string]repo.User{
			"1": {
				ID:   "1",
				Name: "John Reece",
				Roles: map[string]string{
					"1": "user",
					"2": "admin",
				},
			},
			"2": {
				ID:   "2",
				Name: "Bob Smith",
				Roles: map[string]string{
					"1": "admin",
					"2": "user",
				},
			},
		},
		cards,
	)).ServeHTTP(writer, req)

	assert.Equal(t, expectedStatus, writer.Result().StatusCode)
	if expectedStatus != http.StatusOK {
		return
	}

	resBody := TestCardDetails{}
	if err = json.Unmarshal(writer.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}

	assert.Equal(t, cards[reqB.UserID][reqB.CompanyID].MaskedNumber, resBody.MaskedNumber)
	assert.Equal(t, reqB.NewLimit, resBody.Limit)
}

func TestChangeLimitAPI(t *testing.T) {
	testCases := []struct {
		actingUserID   string
		cardUserID     string
		cardCompanyID  string
		newLimit       int
		expectedStatus int
	}{
		{"1", "1", "1", 20000, http.StatusForbidden},
		{"1", "2", "1", 20000, http.StatusForbidden},
		{"1", "1", "2", 20000, http.StatusOK},
		{"1", "2", "2", 20000, http.StatusOK},
		{"2", "1", "1", 20000, http.StatusOK},
		{"2", "2", "1", 20000, http.StatusOK},
		{"2", "1", "2", 20000, http.StatusForbidden},
		{"2", "2", "2", 20000, http.StatusForbidden},
	}
	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("User %s trying to modify limit for user %v and company %v should receive status %v",
				tc.actingUserID, tc.cardUserID, tc.cardCompanyID, tc.expectedStatus),
			func(t *testing.T) {
				testChangeLimitAPI(t, &TestLimitRequest{
					tc.cardUserID,
					tc.cardCompanyID,
					tc.newLimit,
				}, &tc.actingUserID, tc.expectedStatus)
			})
	}
}
