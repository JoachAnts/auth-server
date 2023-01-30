package identity

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/JoachAnts/auth-server/repo"
)

func NewHandler(repo repo.Repo) http.Handler {
	return &h{
		repo: repo,
	}
}

type h struct {
	repo repo.Repo
}

// FIXME make response camel case
func (h *h) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	res := h.repo.GetUser(userID)
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if res.Roles == nil || len(res.Roles) == 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	b, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling identity response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
