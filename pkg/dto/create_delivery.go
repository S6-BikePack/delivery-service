package dto

import (
	"delivery-service/internal/core/domain"
)

type BodyCreateDelivery struct {
	ParcelId      string
	OwnerId       string
	PickupPoint   domain.Location
	DeliveryPoint domain.Location
	PickupTime    int64
}

type ResponseCreateDelivery struct {
	Parcel        domain.Parcel
	Owner         domain.Customer
	PickupPoint   domain.Location
	DeliveryPoint domain.Location
	PickupTime    int64
}

func BuildResponseCreateDelivery(model domain.Delivery) ResponseCreateDelivery {
	return ResponseCreateDelivery{
		Parcel:        model.Parcel,
		Owner:         model.Customer,
		PickupPoint:   model.PickupPoint,
		DeliveryPoint: model.DeliveryPoint,
		PickupTime:    model.PickupTime.Unix(),
	}
}
