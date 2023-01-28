package card

import (
	"encoding/json"
	"log"
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
	// TODO use mux
	if r.Method == http.MethodPost {
		h.doPost(w, r)
		return
	} else if r.Method == http.MethodGet {
		h.doGet(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *h) doGet(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	res := h.repo.GetCard(userID)
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// TODO use camel case
type CardLimitRequest struct {
	UserID   string
	NewLimit int
}

func (h *h) doPost(w http.ResponseWriter, r *http.Request) {
	requestingUserID := r.Header.Get("Authorization")
	if requestingUserID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := h.repo.GetUser(requestingUserID)
	if user.Role != "admin" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var reqBody = CardLimitRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := h.repo.SetCardLimit(reqBody.UserID, reqBody.NewLimit)
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error encoding response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
