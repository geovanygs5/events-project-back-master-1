// Package uc have all the logic related to use cases
package uc

import (
	"LiteraTest/cancel/v1/internal/models"
	"fmt"
	"time"
)

const (
	LayoutFormat    = "2006-01-02 15:04"
	utcTimeZoneName = "UTC"
)

// ParseEventsToUTCUC declaration of use case struct used in this file
type ParseEventsToUTCUC struct{}

// Handle this use case will convert the timezone of each event to UTC to standardize the process
func (uc *ParseEventsToUTCUC) Handle(events models.Events) (models.Events, error) {
	var eventsInUTC models.Events

	for _, event := range events {
		// Original location to get in mind
		location, err := time.LoadLocation(event.Timezone)
		if err != nil {
			return models.Events{}, &models.EventError{
				Code:       models.CodeParseEventError,
				ID:         models.IDDoubleBookedError,
				Message:    fmt.Sprintf("Error setting timezone of event %v", event),
				StatusCode: models.CodeStatusHTTPBusinessError,
			}
		}

		// Start date conversion
		startDateTimeString := event.Start

		originalStartDateTime, err := time.ParseInLocation(LayoutFormat, startDateTimeString, location)
		if err != nil {
			return models.Events{}, &models.EventError{
				Code:       models.CodeParseEventError,
				ID:         models.IDDoubleBookedError,
				Message:    fmt.Sprintf("Error parsing timezone of event %v", event),
				StatusCode: models.CodeStatusHTTPBusinessError,
			}
		}

		startDateTimeInUTC := originalStartDateTime.UTC()
		startDateTimeInUTCString := startDateTimeInUTC.Format(LayoutFormat)

		// End date conversion
		endDateTimeString := event.End

		originalEndDateTime, err := time.ParseInLocation(LayoutFormat, endDateTimeString, location)
		if err != nil {
			return models.Events{}, &models.EventError{
				Code:       models.CodeParseEventError,
				ID:         models.IDDoubleBookedError,
				Message:    fmt.Sprintf("Error parsing timezone of event %v", event),
				StatusCode: models.CodeStatusHTTPBusinessError,
			}
		}

		endDateTimeInUTC := originalEndDateTime.UTC()
		endDateTimeInUTCString := endDateTimeInUTC.Format(LayoutFormat)

		// Append to new list
		eventsInUTC = append(eventsInUTC, models.Event{
			ID:       event.ID,
			Start:    startDateTimeInUTCString,
			End:      endDateTimeInUTCString,
			Timezone: utcTimeZoneName,
		})
	}

	return eventsInUTC, nil
}

// NewParseEventsToUTCUC initialize this use case
func NewParseEventsToUTCUC() *ParseEventsToUTCUC {
	return &ParseEventsToUTCUC{}
}
