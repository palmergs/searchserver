package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"github.com/palmergs/tokensearch"
	"io/ioutil"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	doc := r.Form.Get("q")
	if doc == "" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			doc = string(body)
		}
	}

	pool := tokensearch.NewTokenNodeVisitorPool(root)
	pool.AdvanceThrough(strings.NewReader(doc))
	json.NewEncoder(w).Encode(pool.Matches)
}