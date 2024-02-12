// Package uc have all the logic related to use cases
package uc

import (
	"LiteraTest/create/v1/internal/models"
	"reflect"
	"testing"
)

// TestNewParseEventsToUTCUC test for this method
func TestNewParseEventsToUTCUC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *ParseEventsToUTCUC
	}{
		{
			name: "Success",
			want: NewParseEventsToUTCUC(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewParseEventsToUTCUC(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParseEventsToUTCUC() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseEventsToUTCUC_Handle test for this method
func TestParseEventsToUTCUC_Handle(t *testing.T) {
	t.Parallel()

	type args struct {
		events models.Events
	}

	tests := []struct {
		name    string
		args    args
		want    models.Events
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				events: models.Events{
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
				},
			},
			want: models.Events{
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
			},
			wantErr: false,
		},
		{
			name: "Error loading location",
			args: args{
				events: models.Events{
					models.Event{
						ID:       1,
						Start:    "2023-02-02 13:00",
						End:      "2023-02-02 14:00",
						Timezone: "WRONG",
					},
				},
			},
			want:    models.Events{},
			wantErr: true,
		},
		{
			name: "Error parsing end time",
			args: args{
				events: models.Events{
					models.Event{
						ID:       1,
						Start:    "2023-02-02 14:00",
						End:      "WRONG",
						Timezone: "America/Bogota",
					},
				},
			},
			want:    models.Events{},
			wantErr: true,
		},
		{
			name: "Error parsing start time",
			args: args{
				events: models.Events{
					models.Event{
						ID:       1,
						Start:    "WRONG",
						End:      "2023-02-02 14:00",
						Timezone: "America/Bogota",
					},
				},
			},
			want:    models.Events{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := &ParseEventsToUTCUC{}
			got, err := uc.Handle(tt.args.events)
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
