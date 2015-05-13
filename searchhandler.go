package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"github.com/palmergs/tokensearch"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	doc := r.Form.Get("q")

	pool := tokensearch.NewTokenNodeVisitorPool(root)
	pool.AdvanceThrough(strings.NewReader(doc))
	json.NewEncoder(w).Encode(pool.Matches)
}