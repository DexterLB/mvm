package main

import (
	"net/http"

	"github.com/DexterLB/mvm/imdb/jsonapi"
)

func main() {
	s := &jsonapi.Server{http.DefaultClient}
	http.ListenAndServe(":8088", s)
}
