package main

import (
	"net/http"

	"github.com/acc-event-manager/steam_go"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	opID := steam_go.NewOpenID(r, false)
	switch opID.Mode() {
	case "":
		http.Redirect(w, r, opID.AuthUrl(), http.StatusMovedPermanently)
	case "cancel":
		w.Write([]byte("Authorization cancelled"))
	default:
		steamID, err := opID.ValidateAndGetID()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// Do whatever you want with steam id
		w.Write([]byte(steamID))
	}
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":8081", nil)
}
