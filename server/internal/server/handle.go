package server

import "net/http"

func handlePrefix(mux *http.ServeMux, pattern string, handler http.Handler) {
	mux.Handle(pattern+"/", http.StripPrefix(pattern, handler))
}

func handleFuncExact(mux *http.ServeMux, method string, pattern string, handler http.HandlerFunc) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s := http.StatusMethodNotAllowed
			respondJSON(w, s, newErrorResponse(http.StatusText(s)))
			return
		}

		if r.URL.Path != pattern {
			notFoundHandler(w, r)
			return
		}

		handler(w, r)
	})
}
