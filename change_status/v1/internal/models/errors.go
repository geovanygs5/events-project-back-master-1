// Package internal have all the main logic
package models

// List for code errors related to events
const (
	// CodeFindDoubleBookedError error related to double booked
	CodeFindDoubleBookedError string = "CODE_FIND_DOUBLE_BOOKED_ERROR"
	// CodeParseEventError error related to double booked
	CodeParseEventError string = "CODE_PARSE_EVENT_ERROR"
	// IDDoubleBookedError error related to double booked
	IDDoubleBookedError string = "ID_DOUBLE_BOOKED_ERROR"
	// CodeGeneralError Unexpected errors code
	CodeGeneralError string = "CODE_GENERAL_ERROR"
	// IDGeneralError Unexpected errors ID
	IDGeneralError string = "ID_GENERAL_ERROR"
	// GeneralErrorTitle title for all errors
	GeneralErrorTitle string = "Error"
)

// CodeStatusHTTPBusinessError HTTP Status Code Business Error 280
const CodeStatusHTTPBusinessError int = 280

// GeneralError for unexpected errors
type GeneralError struct {
	Code       string
	ID         string
	Message    string
	StatusCode int
}

// Error get the error message
func (e *GeneralError) Error() string {
	return e.Message
}

// EventError for unexpected errors
type EventError struct {
	Code       string
	ID         string
	Message    string
	StatusCode int
}

// Error get the error message
func (e *EventError) Error() string {
	return e.Message
}

// ErrorJSONAPI struct base from error response
type ErrorJSONAPI struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// ErrorsJSONAPIProvider interface to add or get errors
type ErrorsJSONAPIProvider interface {
	Add(jsonError ErrorJSONAPI) *ErrorsJSONAPI
	Get() *ErrorsJSONAPI
}

// ErrorsJSONAPI list of errors
type ErrorsJSONAPI struct {
	Errors []ErrorJSONAPI `json:"errors"`
}

// Add an error
func (e *ErrorsJSONAPI) Add(jsonError ErrorJSONAPI) *ErrorsJSONAPI {
	e.Errors = append(e.Errors, jsonError)

	return e
}

// Get an error
func (e *ErrorsJSONAPI) Get() *ErrorsJSONAPI {
	return e
}
