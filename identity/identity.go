package identity

import (
	"encoding/json"
	"log"
	"net/http"
)

type IdentityResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func IdentityHandler(w http.ResponseWriter, r *http.Request) {
	res := IdentityResponse{
		ID:   "1",
		Name: "John Smith",
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
