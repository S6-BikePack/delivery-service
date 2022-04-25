package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTimeAndPlace(t *testing.T) {
	type args struct {
		address     string
		coordinates Location
	}
	tests := []struct {
		name    string
		args    args
		want    TimeAndPlace
		wantErr bool
	}{
		{
			name: "PassingTimeAndPlace",
			args: args{
				address: "Test street 1",
				coordinates: Location{
					Latitude:  1,
					Longitude: 2,
				},
			},
			want: TimeAndPlace{
				Coordinates: Location{
					Latitude:  1,
					Longitude: 2,
				},
				Address: "Test street 1",
				Time:    time.Time{},
			},
			wantErr: false,
		},
		{
			name: "MissingAddress",
			args: args{
				coordinates: Location{
					Latitude:  1,
					Longitude: 2,
				},
			},
			want:    TimeAndPlace{},
			wantErr: true,
		},
		{
			name: "MissingCoordinates",
			args: args{
				address: "Test street 1",
			},
			want:    TimeAndPlace{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTimeAndPlace(tt.args.address, tt.args.coordinates)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTimeAndPlace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTimeAndPlace() got = %v, want %v", got, tt.want)
			}
		})
	}
}
