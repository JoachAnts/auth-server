package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JoachAnts/auth-server/card"
	"github.com/JoachAnts/auth-server/identity"
)

func identityHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You")
}

func Start() {
	http.HandleFunc("/me", identity.Handler)
	http.HandleFunc("/card", card.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
