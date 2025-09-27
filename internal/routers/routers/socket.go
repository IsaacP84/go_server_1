package routers

import (
	"net/http"
)

// Index is a custom handler type
type SocketRouter struct{}

// ServeHTTP implements the http.Handler interface
func (h SocketRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ServeFile(w, r, "../static/index.html")
}
