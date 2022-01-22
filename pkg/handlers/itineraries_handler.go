package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	databaase           = "entertainment"
	itineraryCollection = "itineraries"
)

type ItinerariesHandler struct {
	logger      *zap.SugaredLogger
	mongoClient *mongo.Client
}

func NewItinerariesHandler(logger *zap.SugaredLogger, mongoClient *mongo.Client) *ItinerariesHandler {
	return &ItinerariesHandler{
		logger:      logger,
		mongoClient: mongoClient,
	}
}

// CreateItinerary accepts a request and uses the mongo client to create
// a document in MongoDB with the itinerary
//
// swagger:route PUT /itinerary Itinerary CreateItinerary
// Create an itinerary
//
// parameters:
// + in: header
//   name: x-api-key
//   type: string
//   required: true
// responses:
//  '201':
//      description: Itinerary created successfully
//  '401':
//      description: Unauthorized request
func (i *ItinerariesHandler) CreateItinerary(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	itinerary := Itinerary{}

	// Decode the request body into the struct
	// Return 400 if there's an error
	err := json.NewDecoder(r.Body).Decode(&itinerary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(itinerary)

	//Perform InsertOne operation & validate against the error.
	itineraryCollection := i.mongoClient.Database(databaase).Collection(itineraryCollection)
	_, err = itineraryCollection.InsertOne(ctx, itinerary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Inserted itinerary to Mongo")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// GetItinerary accepts a request and searches MongoDB for matching itineraries
// Then returns the itineraries in order of startDate
//
// swagger:route GET /itinerary Itinerary GetItinerary
// Create an itinerary
//
// parameters:
// + in: header
//   name: x-api-key
//   type: string
//   required: true
// responses:
//  200: Itineraries
//  '401':
//      description: Unauthorized request
func (i *ItinerariesHandler) GetItinerary(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	itineraryCollection := i.mongoClient.Database(databaase).Collection(itineraryCollection)

	// Read documents from itinerary collection in Mongo and sort them so earliest start date is first
	opts := options.Find()
	opts.SetSort(bson.D{{"startDate", 1}})
	sortCursor, err := itineraryCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	var itinerariesSorted []bson.M
	if err = sortCursor.All(ctx, &itinerariesSorted); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Retrieved itineraries from Mongo")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(itinerariesSorted)
}
