package routes

import (
	"github.com/gorilla/mux"

	"goshell/internal/handlers"
)

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/", handlers.HandleRoot).Methods("GET", "OPTIONS")
	// r.HandleFunc("/form_types", handlers.HandleFormTypes).Methods("GET", "OPTIONS")
}
