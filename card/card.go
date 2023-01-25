package card

import (
	"encoding/json"
	"net/http"
)

type TestCardResponse struct {
	// TODO think about how to support different currencies
	MaskedNumber string `json:"maskedNumber"`
	Exp          string `json:"exp"`
	Limit        string `json:"limit"`
	Balance      string `json:"balance"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	res := TestCardResponse{
		MaskedNumber: "**** **** **** 4444",
		Exp:          "12/23",
		Limit:        "10000",
		Balance:      "4321",
	}
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
