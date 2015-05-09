package main

import (
	"fmt"
	"flag"
	"net/http"
	"strings"
	"encoding/json"
	"regexp"
	"github.com/palmergs/tokensearch"
)

var root = tokensearch.NewTokenNode()

var validPath = regexp.MustCompile("^/tokens/([a-zA-Z0-9_-]+)$")

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query := r.Form.Get("q")
	fmt.Printf("search :: %v\n", query)

	matches, err := root.Find(query)
	if err != nil {
		http.Error(w, "Not found", 404)
	} else {
		json.NewEncoder(w).Encode(matches)
	}
}

func tokensHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Printf("%v tokens", r.Method)

	switch strings.ToUpper(r.Method) {
	case "POST", "PUT":
		insertTokenHandler(w, r)
	case "GET", "":
		getTokenHandler(w, r)
	}
}

func insertTokenHandler(w http.ResponseWriter, r *http.Request) {
	ident := r.FormValue("ident")
	display := r.FormValue("display")
	category := r.FormValue("category")
	token := tokensearch.NewToken(ident, display, category)
	_, err := root.Insert(token)
	if err != nil {
		http.Error(w, err.Error(), 401)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	ident := validPath.FindStringSubmatch(r.URL.Path)
	if ident == nil {
		http.NotFound(w, r)
	} else {
		// TODO: scan through tree looking for matching strings
		// root.FindTokens(ident)
		fmt.Fprintf(w, "GET Token : %v", ident[1])
	}
}

func main() {

	serverPort := flag.Int("p", 6060, "server port")
	flag.Parse()

	serverAddr := fmt.Sprintf(":%v", *serverPort)
	fmt.Printf("Starting server on %v...\n", serverAddr)

	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/tokens", tokensHandler)
	http.ListenAndServe(serverAddr, nil)

	fmt.Printf("done\n")
}
