package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JoachAnts/auth-server/identity"
)

func identityHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You")
}

func Start() {
	identity.Register()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
