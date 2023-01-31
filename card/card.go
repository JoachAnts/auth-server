package card

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/JoachAnts/auth-server/repo"
)

type CardResponse struct {
	Cards []Card `json:"cards"`
}

type Card struct {
	CompanyID string      `json:"companyID"`
	Card      CardDetails `json:"card"`
}

type CardDetails struct {
	// TODO think about how to support different currencies
	MaskedNumber string `json:"maskedNumber"`
	Exp          string `json:"exp"`
	Limit        int    `json:"limit"`
	Balance      int    `json:"balance"`
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
	res := h.repo.GetCards(userID)
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var cards []Card
	for k, v := range res {
		cards = append(cards, Card{
			CompanyID: k,
			Card: CardDetails{
				MaskedNumber: v.MaskedNumber,
				Exp:          v.Exp,
				Limit:        v.Limit,
				Balance:      v.Balance,
			},
		})
	}
	response := CardResponse{
		Cards: cards,
	}
	b, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// TODO use camel case
type CardLimitRequest struct {
	UserID    string
	CompanyID string
	NewLimit  int
}

func (h *h) doPost(w http.ResponseWriter, r *http.Request) {
	requestingUserID := r.Header.Get("Authorization")
	if requestingUserID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := h.repo.GetUser(requestingUserID)
	var reqBody = CardLimitRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user.Roles[reqBody.CompanyID] != "admin" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	res := h.repo.SetCardLimit(reqBody.UserID, reqBody.CompanyID, reqBody.NewLimit)
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resBody := CardDetails{
		MaskedNumber: res.MaskedNumber,
		Exp:          res.Exp,
		Limit:        res.Limit,
		Balance:      res.Balance,
	}
	b, err := json.Marshal(resBody)
	if err != nil {
		log.Printf("Error encoding response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
