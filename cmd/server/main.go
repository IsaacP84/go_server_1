package server

import (
	"net/http"

	"github.com/isaacp84/go_server_1/internal/routes"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(mux)
	var handler http.Handler = mux

	return handler
}
