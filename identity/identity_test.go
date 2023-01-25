package identity_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JoachAnts/auth-server/identity"
	"github.com/stretchr/testify/assert"
)

func TestIdentity(t *testing.T) {
	writer := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/me", nil)
	if err != nil {
		t.Fatal(err)
	}

	identity.IdentityHandler(writer, req)

	assert.Equal(t, http.StatusOK, writer.Result().StatusCode)
	assert.Equal(t, "{}", writer.Body.String())
}
