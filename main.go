package main

import (
	"entertain/pkg/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	r := mux.NewRouter()

	logger := zap.NewExample().Sugar()
	defer logger.Sync() // flushes buffer, if any

	eventsHandler := handlers.NewEventsHandler(logger)
	restaurantsHandler := handlers.NewRestaurantssHandler(logger)

	// Routes including path and handler
	r.HandleFunc("/events", eventsHandler.GetAllEvents).Methods(http.MethodGet)
	r.HandleFunc("/restaurants", restaurantsHandler.GetRestaurants).Methods(http.MethodGet)

	fmt.Println("Starting server at port 8080")

	// Bind to a port and pass router in
	panic(http.ListenAndServe(":8080", r))
}
