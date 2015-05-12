package main

import (
	"net/http"
	"encoding/json"
	"unicode/utf8"
	"github.com/palmergs/tokensearch"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	doc := r.Form.Get("q")

	allMatches := make([]*tokensearch.TokenMatch, 0)
	onMatch := func(matches []*tokensearch.TokenMatch) {
		if matches != nil && len(matches) > 0 {
			allMatches = append(allMatches, matches...)
		}
	}

	pool := tokensearch.NewTokenNodeVisitorPool(root)
	lastWasChar := true
	charCount := 0
	lastPosition := 0
	for i, w := 0, 0; i < len(doc); i += w {
		runeValue, width := utf8.DecodeRuneInString(doc[i:])
		w = width

		if w > 0 {
			normalizedRune, currIsChar := tokensearch.NormalizeRune(runeValue)
			if currIsChar {

				if charCount > 0 && !lastWasChar {

					// advance for a deferred separator character for existing visitors
					pool.Advance(' ', lastPosition, onMatch)
				}

				if charCount == 0 || !lastWasChar {

					// visitors begin parsing at beginning of valid strings
					pool.InitVisitor(i)
				}

				// advance of token character
				pool.Advance(normalizedRune, i, onMatch)
				charCount++
				lastPosition = i
			}
			lastWasChar = currIsChar
		}
	}

	json.NewEncoder(w).Encode(allMatches)
}