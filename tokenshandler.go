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
	case "DELETE":
		deleteTokenHandler(w, r)
	case "GET", "":
		getTokensHandler(w, r)
	}
}

func deleteTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := unmarshalToken(r)
	if err != nil {
		writeError(w, err)
	} else {
		token.InitKey();
		root.Remove(&token)
		writeToken(w, token)
	}
}

func insertTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := unmarshalToken(r)
	if err != nil {
		writeError(w, err)
	} else {
		token.InitKey();
		root.Insert(&token)
		writeToken(w, token)
	}
}

func getTokensHandler(w http.ResponseWriter, r *http.Request) {
	matches := root.AllValues(9999)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(matches); err != nil {
		panic(err)
	}
}

func unmarshalToken(r *http.Request) (tokensearch.Token, error) {
	var token tokensearch.Token
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4097))
	if err != nil {
		return token, err
	}
	if err := r.Body.Close(); err != nil {
		return token, err
	}
	if err := json.Unmarshal(body, &token); err != nil {
		return token, err
	}
	return token, nil
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(422)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		panic(err)
	}
}

func writeToken(w http.ResponseWriter, token tokensearch.Token) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		panic(err)
	}
}