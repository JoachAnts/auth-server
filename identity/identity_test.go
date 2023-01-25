package identity_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JoachAnts/auth-server/identity"
	"github.com/stretchr/testify/assert"
)

type TestIdentityResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestIdentityValid(t *testing.T) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/me", nil)
	req.Header.Add("authorization", "1")
	if err != nil {
		t.Fatal(err)
	}

	identity.IdentityHandler(writer, req)

	assert.Equal(t, http.StatusOK, writer.Result().StatusCode)
	resBody := TestIdentityResponse{}
	if err = json.Unmarshal(writer.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}
	assert.Equal(t, "1", resBody.ID)
	assert.Equal(t, "John Smith", resBody.Name)
}

func TestIdentityUnauthenticated(t *testing.T) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/me", nil)
	if err != nil {
		t.Fatal(err)
	}

	identity.IdentityHandler(writer, req)

	assert.Equal(t, http.StatusUnauthorized, writer.Result().StatusCode)
}

func TestIdentityNotExisting(t *testing.T) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/me", nil)
	req.Header.Add("authorization", "1001")
	if err != nil {
		t.Fatal(err)
	}

	identity.IdentityHandler(writer, req)

	assert.Equal(t, http.StatusNotFound, writer.Result().StatusCode)
}
