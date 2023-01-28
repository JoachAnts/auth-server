package card

import (
	"encoding/json"
	"net/http"

	"github.com/JoachAnts/auth-server/repo"
)

type CardResponse struct {
	// TODO think about how to support different currencies
	MaskedNumber string `json:"maskedNumber"`
	Exp          string `json:"exp"`
	Limit        string `json:"limit"`
	Balance      string `json:"balance"`
}

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
	res := h.repo.GetCard(userID)
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
