package server

import (
	"net/http"
)

func New() (http.Handler, error) {
	mux := http.NewServeMux()

	// Register routes.
	mux.HandleFunc("/track", trackHandler)

	// Register default 404.
	mux.HandleFunc("/", handleNotFound)

	// Apply middleware.
	h := withRecover(mux)

	return h, nil
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	s := http.StatusNotFound
	respondJSON(w, s, newErrorResponse(http.StatusText(s)))
}

func withRecover(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	}
}
