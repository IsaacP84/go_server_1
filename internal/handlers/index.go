package handlers

import (
	"fmt"
	"net/http"

	"github.com/isaacp84/go_server_1/config"
)

// Index is a custom handler type
type Index struct{}

// ServeHTTP implements the http.Handler interface
func (h Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello from Index! You requested: %s\n", r.URL.Path)
	http.ServeFile(w, r, fmt.Sprintf("%s/index.html", config.LoadedConfig.Server.Static_file_directory))
}
