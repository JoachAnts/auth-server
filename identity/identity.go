package identity

import (
	"encoding/json"
	"log"
	"net/http"
)

type IdentityResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var db = map[string]IdentityResponse{
	"1": {
		ID:   "1",
		Name: "John Smith",
		Role: "user",
	},
	"2": {
		ID:   "2",
		Name: "Bob Bloggs",
		Role: "admin",
	},
}

func IdentityHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
	}
	res, found := db[userID]
	if !found {
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

func Register() {
	http.HandleFunc("/me", IdentityHandler)
}
