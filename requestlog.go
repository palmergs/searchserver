package main

import (
	"log"
	"time"
	"net/http"
)

func RequestLog(inner func(w http.ResponseWriter, r *http.Request), name string) (func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner(w, r)
		log.Printf("%s\t%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				name,
				time.Since(start))
	}
}