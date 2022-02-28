package server

import (
	"net/http"

	"github.com/os4ua/browser-crypto-mining/server/internal/track"
)

func New() (http.Handler, error) {
	tr := track.NewDatastore()

	mux := http.NewServeMux()

	// Register routes.
	handlePrefix(mux, "/track", newTrackHandler(tr))

	// Register default 404.
	mux.HandleFunc("/", notFoundHandler)

	// Apply middleware.
	h := withRecoveryHandler(mux)

	return h, nil
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	s := http.StatusNotFound
	respondJSON(w, s, newErrorResponse(http.StatusText(s)))
}

func withRecoveryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
