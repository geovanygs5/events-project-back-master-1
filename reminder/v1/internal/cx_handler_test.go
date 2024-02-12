// Package internal contains all the main logic
package internal

import (
	"LiteraTest/create/v1/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/mock"
)

// findDoubleBookedEventsUCMock mock for this use case
type findDoubleBookedEventsUCMock struct {
	mock.Mock
}

// Handle mock for this method
func (m *findDoubleBookedEventsUCMock) Handle(events models.Events) (models.DoubleBookedEvents, error) {
	args := m.Called(events)

	return args.Get(0).(models.DoubleBookedEvents), args.Error(1)
}

// parseEventsToUTCUCMock mock for this use case
type parseEventsToUTCUCMock struct {
	mock.Mock
}

// Handle mock for this method
func (m *parseEventsToUTCUCMock) Handle(events models.Events) (models.Events, error) {
	args := m.Called(events)

	return args.Get(0).(models.Events), args.Error(1)
}

// getDataFromGoldenFile This method reads the golden file located in the path given and return the content (string)
func getDataFromGoldenFile(filePath string) string {
	goldenFile, _ := os.Open(filePath)
	defer func(goldenFile *os.File) {
		err := goldenFile.Close()
		if err != nil {
			panic(err)
		}
	}(goldenFile)

	fileBytes, _ := ioutil.ReadAll(goldenFile)
	data := &bytes.Buffer{}
	_ = json.Compact(data, fileBytes)

	return data.String()
}

// TestHandler_Handle Test for this method
func TestHandler_Handle(t *testing.T) {
	type fields struct {
		findDoubleBookedEventsUC *findDoubleBookedEventsUCMock
		parseEventsToUTCUC       *parseEventsToUTCUCMock
	}

	type args struct {
		event events.APIGatewayProxyRequest
	}

	eventsInBogota := models.Events{
		models.Event{
			ID:       1,
			Start:    "2023-02-02 13:00",
			End:      "2023-02-02 14:00",
			Timezone: "America/Bogota",
		},
		models.Event{
			ID:       2,
			Start:    "2023-02-02 16:00",
			End:      "2023-02-02 18:00",
			Timezone: "America/Bogota",
		},
	}

	eventsInUTC := models.Events{
		models.Event{
			ID:       1,
			Start:    "2023-02-02 18:00",
			End:      "2023-02-02 19:00",
			Timezone: "UTC",
		},
		models.Event{
			ID:       2,
			Start:    "2023-02-02 21:00",
			End:      "2023-02-02 23:00",
			Timezone: "UTC",
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
		mock    func(f fields)
	}{
		{
			name: "Success empty response",
			fields: fields{
				findDoubleBookedEventsUC: &findDoubleBookedEventsUCMock{},
				parseEventsToUTCUC:       &parseEventsToUTCUCMock{},
			},
			args: args{
				event: events.APIGatewayProxyRequest{
					Body: getDataFromGoldenFile(
						"./testdata/no_double_booked_request.golden",
					),
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body: getDataFromGoldenFile(
					"./testdata/no_double_booked_response.golden",
				),
			},
			wantErr: false,
			mock: func(f fields) {
				f.parseEventsToUTCUC.On("Handle", eventsInBogota).Once().Return(eventsInUTC, nil)
				f.findDoubleBookedEventsUC.On("Handle", eventsInUTC).Once().
					Return(models.DoubleBookedEvents{}, nil)
			},
		},
		{
			name: "Success with double booked events",
			fields: fields{
				findDoubleBookedEventsUC: &findDoubleBookedEventsUCMock{},
				parseEventsToUTCUC:       &parseEventsToUTCUCMock{},
			},
			args: args{
				event: events.APIGatewayProxyRequest{
					Body: getDataFromGoldenFile(
						"./testdata/no_double_booked_request.golden",
					),
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body: getDataFromGoldenFile(
					"./testdata/no_double_booked_response.golden",
				),
			},
			wantErr: false,
			mock: func(f fields) {
				f.parseEventsToUTCUC.On("Handle", eventsInBogota).Once().Return(eventsInUTC, nil)
				f.findDoubleBookedEventsUC.On("Handle", eventsInUTC).Once().
					Return(models.DoubleBookedEvents{}, nil)
			},
		},
		{
			name: "Fail parse events to utc",
			fields: fields{
				findDoubleBookedEventsUC: &findDoubleBookedEventsUCMock{},
				parseEventsToUTCUC:       &parseEventsToUTCUCMock{},
			},
			args: args{
				event: events.APIGatewayProxyRequest{
					Body: getDataFromGoldenFile(
						"./testdata/request_with_wrong_location.golden",
					),
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: models.CodeStatusHTTPBusinessError,
				Body: getDataFromGoldenFile(
					"./testdata/response_parse_event_error.golden",
				),
			},
			wantErr: false,
			mock: func(f fields) {
				f.parseEventsToUTCUC.On("Handle", models.Events{
					models.Event{
						ID:       1,
						Start:    "2023-02-02 13:00",
						End:      "2023-02-02 14:00",
						Timezone: "WRONG",
					},
				}).Once().Return(models.Events{}, &models.EventError{
					Code: models.CodeParseEventError,
					ID:   models.IDDoubleBookedError,
					Message: fmt.Sprintf("Error setting timezone of event %v",
						models.Event{
							ID:       1,
							Start:    "2023-02-02 13:00",
							End:      "2023-02-02 14:00",
							Timezone: "WRONG",
						}),
					StatusCode: models.CodeStatusHTTPBusinessError,
				})
			},
		},
		{
			name: "Fail by find double booked events uc",
			fields: fields{
				findDoubleBookedEventsUC: &findDoubleBookedEventsUCMock{},
				parseEventsToUTCUC:       &parseEventsToUTCUCMock{},
			},
			args: args{
				event: events.APIGatewayProxyRequest{
					Body: getDataFromGoldenFile(
						"./testdata/no_double_booked_request.golden",
					),
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: models.CodeStatusHTTPBusinessError,
				Body: getDataFromGoldenFile(
					"./testdata/response_find_by_double_booked_error.golden",
				),
			},
			wantErr: false,
			mock: func(f fields) {
				f.parseEventsToUTCUC.On("Handle", eventsInBogota).Once().Return(models.Events{
					models.Event{
						ID:       1,
						Start:    "2023-02-02 18:00",
						End:      "2023-02-02 19:00",
						Timezone: "UTC",
					},
				}, nil)
				f.findDoubleBookedEventsUC.On("Handle", models.Events{
					models.Event{
						ID:       1,
						Start:    "2023-02-02 18:00",
						End:      "2023-02-02 19:00",
						Timezone: "UTC",
					},
				}).Once().
					Return(models.DoubleBookedEvents{}, &models.EventError{
						Code: models.CodeFindDoubleBookedError,
						ID:   models.IDDoubleBookedError,
						Message: fmt.Sprintf("Error parsing timezone of UTC event %v", models.Event{
							ID:       1,
							Start:    "2023-02-02 18:00",
							End:      "2023-02-02 19:00",
							Timezone: "UTC",
						}),
						StatusCode: models.CodeStatusHTTPBusinessError,
					})
			},
		},
		{
			name: "General error response",
			fields: fields{
				findDoubleBookedEventsUC: &findDoubleBookedEventsUCMock{},
				parseEventsToUTCUC:       &parseEventsToUTCUCMock{},
			},
			args: args{
				event: events.APIGatewayProxyRequest{
					Body: getDataFromGoldenFile(
						"./testdata/no_double_booked_request.golden",
					),
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body: getDataFromGoldenFile(
					"./testdata/general_error_response.golden",
				),
			},
			wantErr: true,
			mock: func(f fields) {
				f.parseEventsToUTCUC.On("Handle", eventsInBogota).Once().Return(eventsInUTC, nil)
				f.findDoubleBookedEventsUC.On("Handle", eventsInUTC).Once().
					Return(models.DoubleBookedEvents{}, errors.New("error"))
			},
		},
	}

	for _, tt := range tests {
		tt.mock(tt.fields)
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				findDoubleBookedEventsUC: tt.fields.findDoubleBookedEventsUC,
				parseEventsToUTCUC:       tt.fields.parseEventsToUTCUC,
			}
			got, err := h.Handle(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewHandler Test for this method
func TestNewHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		findDoubleBookedEventsUC FindDoubleBookedEventsUCInterface
		parseEventsToUTCUC       ParseEventsToUTCUCInterface
	}

	arguments := args{
		findDoubleBookedEventsUC: &findDoubleBookedEventsUCMock{},
		parseEventsToUTCUC:       &parseEventsToUTCUCMock{},
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "Success",
			args: arguments,
			want: NewHandler(
				arguments.findDoubleBookedEventsUC,
				arguments.parseEventsToUTCUC,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewHandler(
				tt.args.findDoubleBookedEventsUC,
				tt.args.parseEventsToUTCUC,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
