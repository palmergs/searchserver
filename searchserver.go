package main

import (
	"fmt"
	"net/http"
	"strings"
	"regexp"
	"github.com/palmergs/tokensearch"
)

var root = tokensearch.NewTokenNode()

var validPath = regexp.MustCompile("^/tokens/([a-zA-Z0-9_-]+)$")

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query := r.FormValue("q")
	matches, err := root.Find(query)
	if err != nil {
		http.Error(w, "Not found", 404)
	} else {
		for match := range matches {
			// TODO
			fmt.Printf("Match = %v\n", match)
		}
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	switch strings.ToUpper(r.Method) {
	case "POST", "PUT":
		insertTokenHandler(w, r)
	case "DELETE":
		deleteTokenHandler(w, r)
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
		// TODO
	}
}

func deleteTokenHandler(w http.ResponseWriter, r *http.Request) {
	ident := r.FormValue("ident")
	display := r.FormValue("display")
	token := tokensearch.NewToken(ident, display, "")
	_, err := root.Remove(token)
	if err != nil {
		http.Error(w, err.Error(), 401)
	} else {
		// TODO
	}
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	ident := validPath.FindStringSubmatch(r.URL.Path)
	if ident == nil {
		http.NotFound(w, r)
	} else {
		// TODO: scan through tree looking for matching strings
		// root.FindTokens(ident)
	}
}

func main() {
	fmt.Println("Starting server on port 6060...")
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/tokens/", tokenHandler)
	http.ListenAndServe(":6060", nil)
}