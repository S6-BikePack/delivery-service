package dto

import (
	"delivery-service/internal/core/domain"
)

type BodyAssignRider struct {
	RiderId string
}

type ResponseAssignRider domain.Delivery

func BuildResponseAssignRider(model domain.Delivery) ResponseAssignRider {
	return ResponseAssignRider(model)
}
