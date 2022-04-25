package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestNewDelivery(t *testing.T) {
	deliveryTime := time.Now().Add(time.Hour)

	type args struct {
		parcel      Parcel
		owner       Customer
		pickup      TimeAndPlace
		destination TimeAndPlace
	}
	tests := []struct {
		name    string
		args    args
		want    Delivery
		wantErr bool
	}{
		{
			name: "passingDelivery",
			args: args{
				parcel: Parcel{
					ID: "parcel-1",
				},
				owner: Customer{
					ID: "customer-1",
				},
				pickup: TimeAndPlace{
					Coordinates: Location{
						Latitude:  1,
						Longitude: 2,
					},
					Address: "Test street 1",
					Time:    deliveryTime,
				},
				destination: TimeAndPlace{
					Coordinates: Location{
						Latitude:  3,
						Longitude: 4,
					},
					Address: "Test street 3",
				},
			},
			want: Delivery{
				Parcel: Parcel{
					ID: "parcel-1",
				},
				Customer: Customer{
					ID: "customer-1",
				},
				Pickup: TimeAndPlace{
					Coordinates: Location{
						Latitude:  1,
						Longitude: 2,
					},
					Address: "Test street 1",
					Time:    deliveryTime,
				},
				Destination: TimeAndPlace{
					Coordinates: Location{
						Latitude:  3,
						Longitude: 4,
					},
					Address: "Test street 3",
				},
				Status: 0,
			},
			wantErr: false,
		},
		{
			name: "missingParcel",
			args: args{
				owner: Customer{
					ID: "customer-1",
				},
				pickup: TimeAndPlace{
					Coordinates: Location{
						Latitude:  1,
						Longitude: 2,
					},
					Address: "Test street 1",
					Time:    deliveryTime,
				},
				destination: TimeAndPlace{
					Coordinates: Location{
						Latitude:  3,
						Longitude: 4,
					},
					Address: "Test street 3",
				},
			},
			want:    Delivery{},
			wantErr: true,
		},
		{
			name: "missingDestination",
			args: args{
				parcel: Parcel{
					ID: "parcel-1",
				},
				owner: Customer{
					ID: "customer-1",
				},
				pickup: TimeAndPlace{
					Coordinates: Location{
						Latitude:  1,
						Longitude: 2,
					},
					Address: "Test street 1",
					Time:    deliveryTime,
				},
			},
			want:    Delivery{},
			wantErr: true,
		},
		{
			name: "missingPickup",
			args: args{
				parcel: Parcel{
					ID: "parcel-1",
				},
				owner: Customer{
					ID: "customer-1",
				},
				destination: TimeAndPlace{
					Coordinates: Location{
						Latitude:  3,
						Longitude: 4,
					},
					Address: "Test street 3",
				},
			},
			want:    Delivery{},
			wantErr: true,
		},
		{
			name: "PickupInPast",
			args: args{
				parcel: Parcel{
					ID: "parcel-1",
				},
				owner: Customer{
					ID: "customer-1",
				},
				pickup: TimeAndPlace{
					Coordinates: Location{
						Latitude:  1,
						Longitude: 2,
					},
					Address: "Test street 1",
					Time:    time.Now().Add(-1),
				},
				destination: TimeAndPlace{
					Coordinates: Location{
						Latitude:  3,
						Longitude: 4,
					},
					Address: "Test street 3",
				},
			},
			want:    Delivery{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDelivery(tt.args.parcel, tt.args.owner, tt.args.pickup, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDelivery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDelivery() got = %v, want %v", got, tt.want)
			}
		})
	}
}
