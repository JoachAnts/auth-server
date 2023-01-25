package identity

import (
	"fmt"
	"net/http"
)

func IdentityHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{}")
}

func Register() {
	http.HandleFunc("/me", IdentityHandler)
}
