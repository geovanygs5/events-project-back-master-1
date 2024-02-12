// Package uc have all the logic related to use cases
package uc

import (
	"LiteraTest/create/v1/internal/models"
	"reflect"
	"testing"
)

// TestFindDoubleBookedEventsUC_Handle test for this method
func TestFindDoubleBookedEventsUC_Handle(t *testing.T) {
	type args struct {
		events models.Events
	}

	tests := []struct {
		name    string
		args    args
		want    models.DoubleBookedEvents
		wantErr bool
	}{
		{
			name: "Success with double-booked events",
			args: args{
				models.Events{
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
					models.Event{
						ID:       3,
						Start:    "2023-02-02 20:00",
						End:      "2023-02-02 22:00",
						Timezone: "UTC",
					},
				},
			},
			want:    models.DoubleBookedEvents{{3, 2}},
			wantErr: false,
		},
		{
			name: "Fail by start date time",
			args: args{
				models.Events{
					models.Event{
						ID:       1,
						Start:    "WRONG",
						End:      "2023-02-02 19:00",
						Timezone: "UTC",
					},
					models.Event{
						ID:       2,
						Start:    "2023-02-02 21:00",
						End:      "2023-02-02 23:00",
						Timezone: "UTC",
					},
					models.Event{
						ID:       3,
						Start:    "2023-02-02 20:00",
						End:      "2023-02-02 22:00",
						Timezone: "UTC",
					},
				},
			},
			want:    models.DoubleBookedEvents{{3, 2}},
			wantErr: false,
		},
		{
			name: "Fail by end date time",
			args: args{
				models.Events{
					models.Event{
						ID:       1,
						Start:    "2023-02-02 18:00",
						End:      "WRONG",
						Timezone: "UTC",
					},
					models.Event{
						ID:       2,
						Start:    "2023-02-02 21:00",
						End:      "2023-02-02 23:00",
						Timezone: "UTC",
					},
					models.Event{
						ID:       3,
						Start:    "2023-02-02 20:00",
						End:      "2023-02-02 22:00",
						Timezone: "UTC",
					},
				},
			},
			want:    models.DoubleBookedEvents{{3, 2}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &FindDoubleBookedEventsUC{}
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

// TestNewFindDoubleBookedEventsUC test for this method
func TestNewFindDoubleBookedEventsUC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *FindDoubleBookedEventsUC
	}{
		{
			name: "Success",
			want: NewFindDoubleBookedEventsUC(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewFindDoubleBookedEventsUC(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFindDoubleBookedEventsUC() = %v, want %v", got, tt.want)
			}
		})
	}
}
