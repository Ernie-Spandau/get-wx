package main

import (
	"math"
	"testing"
)

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		wantError bool
	}{
		{
			name:      "valid coordinates",
			latitude:  40.7127,
			longitude: -74.0060,
			wantError: false,
		},
		{
			name:      "latitude out of range",
			latitude:  91.3270,
			longitude: -73.2869,
			wantError: true,
		},
		{
			name:      "longitude out of range",
			latitude:  40.7127,
			longitude: -190.8269,
			wantError: true,
		},
		{
			name:      "latitude too low",
			latitude:  -91,
			longitude: -74.0060,
			wantError: true,
		},
		{
			name:      "north pole boundary",
			latitude:  90,
			longitude: 0,
			wantError: false,
		},
		{
			name:      "south pole boundary",
			latitude:  -90,
			longitude: 0,
			wantError: false,
		},
		{
			name:      "east longitude boundary",
			latitude:  0,
			longitude: 180,
			wantError: false,
		},
		{
			name:      "west longitude boundary",
			latitude:  0,
			longitude: -180,
			wantError: false,
		},
		{
			name:      "latitude NaN",
			latitude:  math.NaN(),
			longitude: 0,
			wantError: true,
		},
		{
			name:      "longitude NaN",
			latitude:  0,
			longitude: math.NaN(),
			wantError: true,
		},
		{
			name:      "latitude infinity",
			latitude:  math.Inf(1),
			longitude: 0,
			wantError: true,
		},
		{
			name:      "longitude infinity",
			latitude:  0,
			longitude: math.Inf(1),
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateCoordinates(test.latitude, test.longitude)
			if test.wantError == true && err == nil || test.wantError == false && err != nil {
				t.Errorf("%s failed", test.name)
			}
		})
	}
}
