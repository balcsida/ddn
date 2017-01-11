package main

import (
	"log"
	"net/http"
	"time"
)

// Logger logs queries to the log with some extra information
func Logger(inner http.Handler, handler string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf("[%s]\t%s\t%s\t%s\t%s\t",
			r.RemoteAddr, r.Method, r.RequestURI, handler, time.Since(start))
	})
}
