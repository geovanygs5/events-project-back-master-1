// Package internal contains all the main logic
package internal

import (
	"LiteraTest/show/v1/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

// Handler declaration of handler struct used in this file
type Handler struct {
	findDoubleBookedEventsUC FindDoubleBookedEventsUCInterface
	parseEventsToUTCUC       ParseEventsToUTCUCInterface
}

// FindDoubleBookedEventsUCInterface interface for this use case
type FindDoubleBookedEventsUCInterface interface {
	Handle(events models.Events) (models.DoubleBookedEvents, error)
}

// ParseEventsToUTCUCInterface interface for this use case
type ParseEventsToUTCUCInterface interface {
	Handle(events models.Events) (models.Events, error)
}

// Handle main method controller to execute this lambda function
func (h *Handler) Handle(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody models.RequestBody

	err := json.Unmarshal([]byte(event.Body), &requestBody)
	if err != nil {
		return responseError(err)
	}

	// Standardize timezone in the events
	eventsInUTC, err := h.parseEventsToUTCUC.Handle(requestBody.Events)
	if err != nil {
		return responseError(err)
	}

	// Get the double booked events
	doubleBookedEvents, err := h.findDoubleBookedEventsUC.Handle(eventsInUTC)
	if err != nil {
		return responseError(err)
	}

	// Prepare and response double booked events
	responseBody := models.ResponseBody{
		DoubleBookedEvents: doubleBookedEvents,
	}

	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		return responseError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}, nil
}

// responseError return response according error type
func responseError(err error) (events.APIGatewayProxyResponse, error) {
	var lambdaError error

	var httpStatusCode int

	errors := new(models.ErrorsJSONAPI)

	switch e := err.(type) {
	case *models.EventError:
		errors.Add(models.ErrorJSONAPI{
			Status: strconv.Itoa(e.StatusCode),
			Code:   e.Code,
			ID:     e.ID,
			Title:  models.GeneralErrorTitle,
			Detail: err.Error(),
		})

		httpStatusCode = e.StatusCode
	default:
		lambdaError = e

		errors.Add(models.ErrorJSONAPI{
			Status: strconv.Itoa(http.StatusInternalServerError),
			Code:   models.CodeGeneralError,
			ID:     models.IDGeneralError,
			Title:  models.GeneralErrorTitle,
			Detail: e.Error(),
		})

		httpStatusCode = http.StatusInternalServerError
	}

	errorsResponse, _ := json.Marshal(errors)

	return events.APIGatewayProxyResponse{
		StatusCode: httpStatusCode,
		Body:       string(errorsResponse),
	}, lambdaError
}

// NewHandler Initialize Handle
func NewHandler(
	findDoubleBookedEventsUC FindDoubleBookedEventsUCInterface,
	parseEventsToUTCUC ParseEventsToUTCUCInterface,
) *Handler {
	return &Handler{
		findDoubleBookedEventsUC: findDoubleBookedEventsUC,
		parseEventsToUTCUC:       parseEventsToUTCUC,
	}
}
