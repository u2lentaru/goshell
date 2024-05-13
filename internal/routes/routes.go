package routes

import (
	"github.com/gorilla/mux"

	"goshell/internal/handlers"
)

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/", handlers.HandleList).Methods("GET")
	r.HandleFunc("/commands", handlers.HandleList).Methods("GET")
	r.HandleFunc("/commands/{id:[0-9]+}", handlers.HandleGetOne).Methods("GET")
	r.HandleFunc("/commands", handlers.HandlePostExec).Methods("POST")
	r.HandleFunc("/cmdrun/{id:[0-9]+}", handlers.HandleExecOne).Methods("GET")
	r.HandleFunc("/cmdrun", handlers.HandleExec).Methods("GET")
	r.HandleFunc("/results", handlers.HandleResults).Methods("GET")
}
