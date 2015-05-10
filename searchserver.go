package main

import (
	"fmt"
	"flag"
	"net/http"
	"strings"
	"encoding/json"
	"regexp"
	"github.com/palmergs/tokensearch"
	"io"
	"io/ioutil"
	"unicode/utf8"
)

var root = tokensearch.NewTokenNode()

var validPath = regexp.MustCompile("^/tokens/([a-zA-Z0-9_-]+)$")

type TokenMatch struct {
	Matches 	[]*tokensearch.Token
	StartPos	int
	EndPos		int
}


func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query := r.Form.Get("q")
	fmt.Printf("search :: %v\n", query)

	allMatches := make([]*TokenMatch, 0)
	onMatch := func(matches []*tokensearch.Token, startPos int, endPos int) {
		if matches != nil && len(matches) > 0 {
			tokenMatch := &TokenMatch{Matches: matches, StartPos: startPos, EndPos: endPos}
			allMatches = append(allMatches, tokenMatch)
		}
	}

	pool := tokensearch.NewTokenNodeVisitorPool(root)
	for i, w := 0, 0; i < len(query); i += w {
		runeValue, width := utf8.DecodeRuneInString(query[i:])
		w = width

		pool.Advance(runeValue, i, onMatch)
	}

	json.NewEncoder(w).Encode(allMatches)
}

func tokensHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Printf("%v tokens", r.Method)

	switch strings.ToUpper(r.Method) {
	case "POST", "PUT":
		insertTokenHandler(w, r)
	case "GET", "":
		getTokensHandler(w, r)
	}
}

func insertTokenHandler(w http.ResponseWriter, r *http.Request) {

	var token tokensearch.Token
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4097))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &token); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	root.Insert(&token)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		panic(err)
	}
}

func getTokensHandler(w http.ResponseWriter, r *http.Request) {
	matches := root.AllValues(999)
	json.NewEncoder(w).Encode(matches)
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
