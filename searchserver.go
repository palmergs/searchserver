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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	doc := r.Form.Get("q")
	fmt.Printf("search :: %v\n", doc)

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

func tokensHandler(w http.ResponseWriter, r *http.Request) {

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

func loadFile(pathToFile string) {

	body, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		panic(err)
	}

	count := 0
	var importJson []interface{}
	json.Unmarshal(body, &importJson)

	for _, mapJson := range importJson {
		m := mapJson.(map[string]interface{})
		_, err := root.Insert(tokensearch.NewToken(
				fmt.Sprintf("%.f", m["id"].(float64)),
				m["label"].(string),
				m["category"].(string)))
		if err == nil {
			count++
		}
	}
	fmt.Printf("Inserted %d values\n", count)
}

func getTokensHandler(w http.ResponseWriter, r *http.Request) {
	matches := root.AllValues(999)
	json.NewEncoder(w).Encode(matches)
}

func main() {

	serverPort := flag.Int("p", 6060, "server port")
	importFile := flag.String("f", "", "prepopulate with file")
	flag.Parse()

	serverAddr := fmt.Sprintf(":%v", *serverPort)
	fmt.Printf("Starting server on %v...\n", serverAddr)

	if *importFile != "" {
		fmt.Printf("Prepopulate tree with %s...\n", *importFile)
		loadFile(*importFile)
	}

	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/tokens", tokensHandler)
	http.ListenAndServe(serverAddr, nil)

	fmt.Printf("done\n")
}
