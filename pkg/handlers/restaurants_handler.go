package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

const (
	yelpPath           = "https://api.yelp.com/v3"
	yelpBusinessesPath = yelpPath + "/businesses/search"
)

type RestaurantsHandler struct {
	logger *zap.SugaredLogger
}

func NewRestaurantssHandler(logger *zap.SugaredLogger) *EventsHandler {
	return &EventsHandler{
		logger: logger,
	}
}

type RestaurantsResponse struct {
	Businesses []Business `json:"businesses"`
}

// swagger:model Business
type Business struct {
	Name        string  `json:"name"`
	Url         string  `json:"url"`
	ReviewCount int     `json:"review_count"`
	Rating      float32 `json:"rating"`
	Price       string  `json:"price"`
	Phone       string  `json:"phone"`
	Distance    float64 `json:"distance"`
}

// GetRestuarants accepts a request and calls the Yelp api to get restuarants
//
// swagger:route GET /restaurants Restaurants GetRestuarants
// Get all restaurants
//
// parameters:
// + in: header
//   name: x-api-key
//   type: string
//   required: true
// responses:
//  200: RestaurantsResponse
//  '401':
//      description: Unauthorized request
func (e *EventsHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}

	apiKey := r.Header.Get("X-Api-Key")
	if apiKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	queryParams := r.URL.Query()
	longitude := queryParams.Get("longitude")
	latitude := queryParams.Get("latitude")
	radius := queryParams.Get("radius")
	limit := queryParams.Get("limit")
	var url string

	// Create request for Yelp. If a limit is given, sort by the top rated and use that limit to return the top n rated resturants
	if limit == "" {
		url = fmt.Sprintf("%s?term=restaurants&longitude=%s&latitude=%s&radius=%s", yelpBusinessesPath, longitude, latitude, radius)
	} else {
		url = fmt.Sprintf("%s?term=restaurants&longitude=%s&latitude=%s&radius=%s&sort_by=rating&limit=%s", yelpBusinessesPath, longitude, latitude, radius, limit)
	}

	// Build request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Add auth header
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	} else if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response to struct
	var restaurants RestaurantsResponse
	json.Unmarshal(bodyBytes, &restaurants)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurants)
}
