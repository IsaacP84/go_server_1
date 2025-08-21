package routers

import (
	"net/http"
)

// MyHandler is a custom handler type
type Root struct{}

// ServeHTTP implements the http.Handler interface
func (h Root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
