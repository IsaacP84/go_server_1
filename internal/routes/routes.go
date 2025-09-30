package routes

import (
	"net/http"

	"github.com/isaacp84/go_server_1/config"
	"github.com/isaacp84/go_server_1/internal/handlers"
)

func AddRoutes(
	mux *http.ServeMux,
) {
	mux.Handle("/", &handlers.Index{})

	mux.Handle("/about/", http.NotFoundHandler())

	fs := http.FileServer(http.Dir(config.LoadedConfig.Server.Static_file_directory))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

}
