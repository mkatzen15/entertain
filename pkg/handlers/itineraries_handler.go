package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
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
	// 	StartDate: "1-8-2022",
	// 	DayPlans: []DayPlan{
	// 		{
	// 			Time: "12421",
	// 			Location: EventLocation{
	// 				Name: "test location",
	// 				Url:  "google.com",
	// 			},
	// 		},
	// 		{
	// 			Time: "32532",
	// 			Location: EventLocation{
	// 				Name: "test location2",
	// 				Url:  "google.com/2",
	// 			},
	// 		},
	// 	},
	// }
	fmt.Println(r.Body)

	// Decode the request body into the struct
	// Return 400 if there's an error
	err := json.NewDecoder(r.Body).Decode(&itinerary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(itinerary)

	//Perform InsertOne operation & validate against the error.
	collection := i.mongoClient.Database(databaase).Collection(itineraryCollection)
	_, err = collection.InsertOne(ctx, itinerary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Inserted doc to Mongo")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
