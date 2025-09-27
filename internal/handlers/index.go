package handlers

import (
	"net/http"
)

// Index is a custom handler type
type Index struct{}

// ServeHTTP implements the http.Handler interface
func (h Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ServeFile(w, r, "web/static/index.html")
}
