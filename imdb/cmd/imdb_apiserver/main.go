package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DexterLB/mvm/imdb/jsonapi"
)

func serve(address string) error {
	s := &jsonapi.Server{http.DefaultClient}
	return http.ListenAndServe(address, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: $0 <bind address>\n")
		os.Exit(2)
	}
	err := serve(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
