// Package uc have all the logic related to use cases
package uc

import (
	"LiteraTest/cancel/v1/internal/models"
	"fmt"
	"sync"
	"time"
)

// FindDoubleBookedEventsUC declaration of use case struct used in this file
type FindDoubleBookedEventsUC struct{}

// Handle find all the double booked events in the list of events given
func (uc *FindDoubleBookedEventsUC) Handle(events models.Events) (models.DoubleBookedEvents, error) {
	// variables used for manage concurrency
	var (
		waitGroup sync.WaitGroup
		mutex     sync.Mutex
	)

	waitGroup.Add(len(events))

	var doubleBookedEvents models.DoubleBookedEvents

	// We will loop the events, and we compare each event if overlapping another event
	for _, event := range events {
		// We will create one go routine for each event in the list to check overlapping
		go func(event models.Event) {
			defer waitGroup.Done()

			for _, eventToCheck := range events {
				if event.ID == eventToCheck.ID {
					continue
				}

				// We check if the event is a double booking adding in a list if its match
				isDoubleBookedCheck, _ := inTimeSpan(event, eventToCheck)
				if isDoubleBookedCheck {
					mutex.Lock()
					if !isAlreadyInList(event.ID, eventToCheck.ID, doubleBookedEvents) {
						doubleBookedEvents = append(doubleBookedEvents, []int{event.ID, eventToCheck.ID})
					}
					mutex.Unlock()
				}
			}
		}(event)
	}

	waitGroup.Wait()

	return doubleBookedEvents, nil
}

// isAlreadyInList check if an events pair is already in the list of events given
func isAlreadyInList(eventID, eventToCheckID int, eventsList models.DoubleBookedEvents) bool {
	for _, pair := range eventsList {
		if (eventID == pair[0] && eventToCheckID == pair[1]) || (eventID == pair[1] && eventToCheckID == pair[0]) {
			return true
		}
	}

	return false
}

// inTimeSpan check if an event is overlapping another event
func inTimeSpan(event, eventToCheck models.Event) (bool, error) {
	startToCheck, err := time.Parse(LayoutFormat, eventToCheck.Start)
	if err != nil {
		return false, &models.EventError{
			Code:       models.CodeFindDoubleBookedError,
			ID:         models.IDDoubleBookedError,
			Message:    fmt.Sprintf("Error parsing timezone of UTC event %v", event),
			StatusCode: models.CodeStatusHTTPBusinessError,
		}
	}

	endToCheck, err := time.Parse(LayoutFormat, eventToCheck.End)
	if err != nil {
		return false, &models.EventError{
			Code:       models.CodeFindDoubleBookedError,
			ID:         models.IDDoubleBookedError,
			Message:    fmt.Sprintf("Error parsing timezone of UTC event %v", event),
			StatusCode: models.CodeStatusHTTPBusinessError,
		}
	}

	start, err := time.Parse(LayoutFormat, event.Start)
	if err != nil {
		return false, &models.EventError{
			Code:       models.CodeFindDoubleBookedError,
			ID:         models.IDDoubleBookedError,
			Message:    fmt.Sprintf("Error parsing timezone of UTC event %v", event),
			StatusCode: models.CodeStatusHTTPBusinessError,
		}
	}

	end, err := time.Parse(LayoutFormat, event.End)
	if err != nil {
		return false, &models.EventError{
			Code:       models.CodeFindDoubleBookedError,
			ID:         models.IDDoubleBookedError,
			Message:    fmt.Sprintf("Error parsing timezone of UTC event %v", event),
			StatusCode: models.CodeStatusHTTPBusinessError,
		}
	}

	if (start.After(startToCheck) && start.Before(endToCheck)) ||
		(end.After(startToCheck) && end.Before(endToCheck)) {
		return true, nil
	}

	return false, nil
}

// NewFindDoubleBookedEventsUC initialize this use case
func NewFindDoubleBookedEventsUC() *FindDoubleBookedEventsUC {
	return &FindDoubleBookedEventsUC{}
}
