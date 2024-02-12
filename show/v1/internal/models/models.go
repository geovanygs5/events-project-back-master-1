// Package internal have all the main logic
package models

// RequestBody struct for request body
type RequestBody struct {
	Events Events `json:"events"`
}

// Events declare a list of events
type Events []Event

// Event declare structure for each event
type Event struct {
	ID       int    `json:"id"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Timezone string `json:"timezone"`
}

// DoubleBookedEvents declare a list of pairs of double-booked events
type DoubleBookedEvents [][]int

// ResponseBody struct for response body
type ResponseBody struct {
	DoubleBookedEvents DoubleBookedEvents `json:"double_booked_events"`
}
