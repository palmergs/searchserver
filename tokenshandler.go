package main

import (
	"net/http"
	"strings"
	"io"
	"io/ioutil"
	"encoding/json"
	"github.com/palmergs/tokensearch"
)

func TokensHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
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

	token.InitKey();
	root.Insert(&token)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		panic(err)
	}
}

func getTokensHandler(w http.ResponseWriter, r *http.Request) {
	matches := root.AllValues(9999)
	json.NewEncoder(w).Encode(matches)
}