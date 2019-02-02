package main

import (
	"net/http"
	"os"
	"strings"
)

func basicCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("access-control-allow-origin", "*")

		next.ServeHTTP(w, r)
	})
}

func whitelistHostname(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow := false
		allowedHosts := strings.Split(os.Getenv("ALLOWED_HOSTS"), "#")

		for _, host := range allowedHosts {
			if host == r.Host {
				allow = true
				break
			}
		}

		if !allow {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
