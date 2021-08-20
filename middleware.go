package main

import "net/http"

// CORS handler, sets the context response access-control headers.
//
// Used for all API calls to allow access from all origins.
func HandleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, session")
		writer.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, DELETE, OPTIONS")
		writer.Header().Set("Access-Control-Max-Age", "600")
		next.ServeHTTP(writer, request)
	})
}

// Handle options requests by responding OK to all of them.
//
// Avoids OPTIONS getting into the access logs.
func HandleOptions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "OPTIONS" {
			writer.WriteHeader(http.StatusOK)
		} else {
			next.ServeHTTP(writer, request)
		}
	})
}
