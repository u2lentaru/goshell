package routes

import (
	"github.com/gorilla/mux"

	"goshell/internal/handlers"
)

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/", handlers.HandlerList).Methods("GET")
	r.HandleFunc("/commands", handlers.HandlerList).Methods("GET")
	r.HandleFunc("/commands/{id:[0-9]+}", handlers.HandlerGetOne).Methods("GET")
	r.HandleFunc("/commands", handlers.HandlerPostExec).Methods("POST")
	r.HandleFunc("/cmdrun/{id:[0-9]+}", handlers.HandlerExecOne).Methods("GET")
	r.HandleFunc("/cmdrun", handlers.HandlerExec).Methods("GET")
	r.HandleFunc("/results", handlers.HandlerResults).Methods("GET")
	r.HandleFunc("/health-check", handlers.HealthCheckHandler).Methods("GET")
}
