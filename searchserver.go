package main

import (
	"fmt"
	"log"
	"flag"
	"net/http"
	"regexp"
	"github.com/palmergs/tokensearch"
)

var root = tokensearch.NewTokenNode()

var validPath = regexp.MustCompile("^/tokens/([a-zA-Z0-9_-]+)$")

func main() {

	serverPort := flag.Int("p", 6060, "server port")
	importFile := flag.String("f", "", "prepopulate with file")
	flag.Parse()

	serverAddr := fmt.Sprintf(":%v", *serverPort)
	log.Printf("Starting server on %v...\n", serverAddr)

	if *importFile != "" {
		log.Printf("Prepopulate tree with %s...\n", *importFile)
		root.InsertFromFile(*importFile)
	}

	http.HandleFunc("/search", RequestLog(SearchHandler, "search"))
	http.HandleFunc("/tokens", RequestLog(TokensHandler, "tokens"))
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
