package main

import (
	"context"
	"entertain/pkg/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mongo struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
}

func main() {
	// Get config from config.yaml
	configFile, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	var config Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Setup Mongo client
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.km2hn.mongodb.net/myFirstDatabase?retryWrites=true&w=majority", config.Mongo.Username, config.Mongo.Password)
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)

	// List Mongo dbs for debugging
	databases, err := mongoClient.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	r := mux.NewRouter()

	logger := zap.NewExample().Sugar()
	defer logger.Sync() // flushes buffer, if any

	eventsHandler := handlers.NewEventsHandler(logger)
	restaurantsHandler := handlers.NewRestaurantssHandler(logger)
	itinerariesHandler := handlers.NewItinerariesHandler(logger, mongoClient)

	// Routes including path and handler
	r.HandleFunc("/events", eventsHandler.GetAllEvents).Methods(http.MethodGet)
	r.HandleFunc("/restaurants", restaurantsHandler.GetRestaurants).Methods(http.MethodGet)
	r.HandleFunc("/itinerary", itinerariesHandler.CreateItinerary).Methods(http.MethodPut)
	r.HandleFunc("/itinerary", itinerariesHandler.GetItinerary).Methods(http.MethodGet)

	fmt.Println("Starting server at port 8080")

	// Bind to a port and pass router in
	panic(http.ListenAndServe(":8080", r))
}
