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
			Role: "user",
			Roles: []repo.Role{
				{
					CompanyID: "1",
					Role:      "user",
				},
				{
					CompanyID: "2",
					Role:      "admin",
				},
			},
		},
		"2": {
			ID:   "2",
			Name: "Bob Bloggs",
			Role: "admin",
			Roles: []repo.Role{
				{
					CompanyID: "1",
					Role:      "admin",
				},
				{
					CompanyID: "2",
					Role:      "user",
				},
			},
		},
	}, map[string]repo.Card{
		"1": {
			MaskedNumber: "**** **** **** 4444",
			Exp:          "12/23",
			Limit:        10000,
			Balance:      4321,
		},
	})
	identityHandler := identity.NewHandler(repo)
	cardHandler := card.NewHandler(repo)
	http.HandleFunc("/me", identityHandler.ServeHTTP)
	http.HandleFunc("/card", cardHandler.ServeHTTP)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
