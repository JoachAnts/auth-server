package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JoachAnts/auth-server/card"
	"github.com/JoachAnts/auth-server/identity"
	"github.com/JoachAnts/auth-server/repo"
)

func identityHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You")
}

func Start() {
	identityHandler := identity.NewHandler(repo.NewRepo(map[string]repo.User{
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
	}))
	http.HandleFunc("/me", identityHandler.ServeHTTP)
	http.HandleFunc("/card", card.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
