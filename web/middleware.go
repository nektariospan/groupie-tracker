package web

import (
	"log"
	"net/http"
	"runtime/debug"
)

// RecoveryMiddleware recovers from panics and returns a 500 page
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("🔥 Panic recovered: %v\n%s", err, debug.Stack())
				renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
