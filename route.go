package main

import "net/http"

// Structure to declare a route of the application
type Route struct {
	// Type of the API route.
	Type string

	// Path of the API.
	Path string

	// Handler to process the API call
	Handler http.HandlerFunc
}
