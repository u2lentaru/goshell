package routes

import (
	"github.com/gorilla/mux"

	"goshell/internal/handlers"
)

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/", handlers.HandleRoot).Methods("GET", "OPTIONS")
	r.HandleFunc("/command", handlers.HandleExec).Methods("GET", "OPTIONS")
	r.HandleFunc("/commands", handlers.HandleList).Methods("GET", "OPTIONS")
	r.HandleFunc("/commands/{id:[0-9]+}", handlers.HandleGetOne).Methods("GET", "OPTIONS")
}
