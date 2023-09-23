# steam_go

Simple Steam auth util in Golang. Forked and better version.

### Installation
```
$ go get github.com/pektezol/steam_go
```
### Usage
Just <code>go run main.go</code> in example dir and open [localhost:8081/login](http://localhost:8081/login) link to see how it works

Code from ./example/main.go:
```
package main

import (
	"net/http"

	"github.com/pektezol/steam_go"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	opID := steam_go.NewOpenID(r)
	switch opID.Mode() {
	case "":
		http.Redirect(w, r, opID.AuthUrl(), 301)
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
```
