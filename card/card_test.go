package card_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JoachAnts/auth-server/card"
	"github.com/stretchr/testify/assert"
)

type TestCardResponse struct {
	// TODO think about how to support different currencies
	MaskedNumber string `json:"maskedNumber"`
	Exp          string `json:"exp"`
	Limit        string `json:"limit"`
	Balance      string `json:"balance"`
}

func TestCardUser(t *testing.T) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/card", nil)
	req.Header.Add("authorization", "1")
	if err != nil {
		t.Fatal(err)
	}

	card.Handler(writer, req)

	assert.Equal(t, http.StatusOK, writer.Result().StatusCode)
	resBody := TestCardResponse{}
	if err = json.Unmarshal(writer.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}
	assert.Equal(t, "**** **** **** 4444", resBody.MaskedNumber)
	assert.Equal(t, "12/23", resBody.Exp)
	assert.Equal(t, "10000", resBody.Limit)
	assert.Equal(t, "4321", resBody.Balance)
}
