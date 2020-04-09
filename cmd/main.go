package main

import (
	"flag"
	"net/http"

	gsd "github.com/ScentWoman/GSD"
)

var (
	credentials = flag.String("cred", "credentials.json", "Your credential file from Google API console.")
	token       = flag.String("token", "token.json", "Your app token, if not exist, the program will generate from CLI.")
	listen      = flag.String("listen", "127.0.0.1:8080", "listen address:port")
)

func main() {
	flag.Parse()

	gsd.HandleFunc("/", *credentials, *token)
	http.ListenAndServe(*listen, nil)
}
