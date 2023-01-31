package server

import (
	"log"
	"net/http"

	"github.com/JoachAnts/auth-server/card"
	"github.com/JoachAnts/auth-server/identity"
	"github.com/JoachAnts/auth-server/repo"
)

func Start() {
	repo := repo.NewRepo(map[string]repo.User{
		"1": {
			ID:   "1",
			Name: "John Smith",
			Roles: map[string]string{
				"1": "user",
				"2": "admin",
			},
		},
		"2": {
			ID:   "2",
			Name: "Bob Bloggs",
			Roles: map[string]string{
				"1": "admin",
				"2": "user",
			},
		},
	}, map[string](map[string]repo.Card){
		"1": {
			"1": {
				MaskedNumber: "**** **** **** 1111",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
			"2": {
				MaskedNumber: "**** **** **** 2222",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
		},
		"2": {
			"1": {
				MaskedNumber: "**** **** **** 3333",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
			"2": {
				MaskedNumber: "**** **** **** 4444",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
		},
	})
	identityHandler := identity.NewHandler(repo)
	cardHandler := card.NewHandler(repo)
	http.HandleFunc("/me", identityHandler.ServeHTTP)
	http.HandleFunc("/card", cardHandler.ServeHTTP)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
