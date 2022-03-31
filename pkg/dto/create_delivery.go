package dto

import (
	"delivery-service/internal/core/domain"
)

type BodyCreateDelivery struct {
	ParcelId      string
	PickupPoint   domain.Location
	DeliveryPoint domain.Location
	PickupTime    int64
}

type ResponseCreateDelivery domain.Delivery

func BuildResponseCreateDelivery(model domain.Delivery) ResponseCreateDelivery {
	return ResponseCreateDelivery(model)
}
