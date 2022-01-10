package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

const (
	tmPath      = "https://app.ticketmaster.com/discovery/v2"
	tmEventPath = tmPath + "/events"
)

type EventsHandler struct {
	logger *zap.SugaredLogger
}

func NewEventsHandler(logger *zap.SugaredLogger) *EventsHandler {
	return &EventsHandler{
		logger: logger,
	}
}

type EventResponse struct {
	Embedded EventList `json:"_embedded"`
}

// swagger:model EventList
type EventList struct {
	Events []Event `json:"events"`
}

type Event struct {
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Dates    Date     `json:"dates"`
	Embedded Embedded `json:"_embedded"`
}

type Date struct {
	Start    StartData `json:"start"`
	Timezone string    `json:"timezone"`
	Status   Status    `json:"status"`
}

type StartData struct {
	LocalDate      string `json:"localDate"`
	DateTBD        bool   `json:"dateTBD"`
	TimeTBA        bool   `json:"timeTBA"`
	NoSpecificTime bool   `json:"noSpecificTime"`
}

type Status struct {
	Code string `json:"code"`
}

type Embedded struct {
	Venues []Venue `json:"venues"`
}

type Venue struct {
	Name     string   `json:"name"`
	Url      string   `json:"url"`
	City     Name     `json:"city"`
	State    Name     `json:"state"`
	Address  Address  `json:"address"`
	Location Location `json:"location"`
}

type Name struct {
	Name string `json:"name"`
}

type Address struct {
	Line1 string `json:"line1"`
}

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// GetAllEvents accepts a request and calls the Ticketmaster api to get all events
//
// swagger:route GET /events Events GetEvents
// Get all events
//
// parameters:
// + in: header
//   name: x-api-key
//   type: string
//   required: true
// responses:
//  200: EventList
//  '401':
//      description: Unauthorized request
func (e *EventsHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-Api-Key")
	if apiKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	queryParams := r.URL.Query()
	city := queryParams.Get("city")

	// Call TM api with apiKey
	resp, err := http.Get(fmt.Sprintf("%s?apikey=%s&city=%s", tmEventPath, apiKey, city))
	if err != nil {
		fmt.Println(err.Error())
	} else if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response to struct
	var events EventResponse
	json.Unmarshal(bodyBytes, &events)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events.Embedded)
}
