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

func (h *h) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
	}
	res := h.repo.GetUser(userID)
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
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
