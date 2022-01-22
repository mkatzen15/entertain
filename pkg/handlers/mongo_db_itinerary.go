package handlers

type Itinerary struct {
	StartDate string    `json:"startDate" bson:"startDate"`
	DayPlans  []DayPlan `json:"dayPlans" bson:"dayPlans"`
}

type DayPlan struct {
	Time     string        `json:"time" bson:"time"`
	Location EventLocation `json:"location" bson:"location"`
}

type EventLocation struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}
