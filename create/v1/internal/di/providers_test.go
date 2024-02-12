// Package di have all the injections dependency logic
package di

import (
	"reflect"
	"testing"
)

// Test_newAWSSessionProvider test for aws session provider.
func Test_newAWSSessionProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want SessionProvider
	}{
		{
			name: "ok",
			want: newSessionProvider(&SessionConfig{}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := newAWSSessionProvider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAWSSessionProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
