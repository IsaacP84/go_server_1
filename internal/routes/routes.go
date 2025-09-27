package routes

import (
	"net/http"

	"github.com/isaacp84/go_server_1/internal/handlers"
)

func AddRoutes(
	mux *http.ServeMux,
) {
	mux.Handle("/about", http.NotFoundHandler())
	mux.Handle("/", handlers.Index{})

}
