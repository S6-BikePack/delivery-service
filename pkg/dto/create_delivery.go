package dto

import (
	"delivery-service/internal/core/domain"
)

type BodyCreateDeliveryPickup struct {
	Coordinates domain.Location `json:"coordinates"`
	Address     string          `json:"address"`
	Time        int64           `json:"time"`
}

type BodyCreateDeliveryDestination struct {
	Coordinates domain.Location `json:"coordinates"`
	Address     string          `json:"address"`
}

type BodyCreateDelivery struct {
	ParcelId    string                        `json:"parcelId"`
	OwnerId     string                        `json:"ownerId"`
	Pickup      BodyCreateDeliveryPickup      `json:"pickup"`
	Destination BodyCreateDeliveryDestination `json:"destination"`
}

type ResponseCreateDelivery struct {
	Parcel      domain.Parcel       `json:"parcel"`
	Owner       domain.Customer     `json:"owner"`
	Pickup      domain.TimeAndPlace `json:"pickup"`
	Destination domain.TimeAndPlace `json:"destination"`
}

func BuildResponseCreateDelivery(model domain.Delivery) ResponseCreateDelivery {
	return ResponseCreateDelivery{
		Parcel:      model.Parcel,
		Owner:       model.Customer,
		Pickup:      model.Pickup,
		Destination: model.Destination,
	}
}
