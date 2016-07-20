package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DexterLB/mvm/imdb"
)

// Server implements a JSON API for querying IMDB
type Server struct {
	ImdbClient imdb.HttpGetter
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/item", s.item)
	mux.HandleFunc("/search", s.search)

	mux.ServeHTTP(w, r)
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.searchGet(w, r)
	case "POST":
		s.searchPost(w, r)
	default:
		http.Error(
			w,
			fmt.Sprintf("wrong method: %s", r.Method),
			500,
		)
	}
}

func (s *Server) searchGet(w http.ResponseWriter, r *http.Request) {
	var err error

	parameters := r.URL.Query()
	query := &imdb.SearchQuery{}

	query.Query = parameters.Get("query")

	if len(parameters.Get("year")) != 0 {
		query.Year, err = strconv.Atoi(parameters.Get("year"))
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("year parameter missing or not integer: %s"),
				500,
			)
			return
		}
	}

	if len(parameters.Get("category")) != 0 {
		err = json.Unmarshal([]byte(parameters.Get("category")), &query.Category)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("unable to parse category parameter: %s", err),
				500,
			)
			return
		}
	}

	if len(parameters.Get("exact")) != 0 {
		err = json.Unmarshal([]byte(parameters.Get("exact")), &query.Exact)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("unable to parse exact parameter: %s", err),
				500,
			)
			return
		}
	}

	data, err := imdb.SearchWithClient(query, s.ImdbClient)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("unable to get data from imdb: %s", err),
			500,
		)
		return
	}

	writeData(w, data)
}

func (s *Server) searchPost(w http.ResponseWriter, r *http.Request) {
	query := &imdb.SearchQuery{}
	err := json.NewDecoder(r.Body).Decode(query)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("unable to parse request: %s", err),
			500,
		)
		return
	}

	data, err := imdb.SearchWithClient(query, s.ImdbClient)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("unable to get data from imdb: %s", err),
			500,
		)
		return
	}

	writeData(w, data)
}

func (s *Server) item(w http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()
	id, err := strconv.Atoi(parameters.Get("id"))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("ID parameter missing or not integer: %s", err),
			500,
		)
		return
	}

	data, err := imdb.NewWithClient(id, s.ImdbClient).AllData()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Unable to get data from imdb: %s", err),
			500,
		)
		return
	}

	writeData(w, data)
}

func writeData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Unable to encode data from imdb: %s", err),
			500,
		)
		return
	}
}
