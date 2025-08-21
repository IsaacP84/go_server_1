package routers

import (
	"fmt"
	"html"
	"net/http"
)

// MyHandler is a custom handler type
type Admin struct{}

// ServeHTTP implements the http.Handler interface
func (h Admin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from MyHandler! You requested: %s", r.URL.Path)
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
