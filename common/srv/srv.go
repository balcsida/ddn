package srv

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// Logger logs queries to the log with some extra information
func Logger(inner http.Handler, handler string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		if strings.HasPrefix(r.RequestURI, "/alive") ||
			r.RequestURI == "/heartbeat" {
			return
		}

		log.Printf("[%s]\t%s\t%s\t%s\t",
			r.RemoteAddr, r.Method, r.RequestURI, time.Since(start))
	})
}
