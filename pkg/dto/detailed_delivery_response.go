package dto

import (
	"delivery-service/internal/core/domain"
	"time"
)

type DetailedDeliveryResponseTimeAndPlace struct {
	Coordinates domain.Location `json:"coordinates"`
	Address     string          `json:"address"`
	Time        time.Time       `json:"time"`
}

type DetailedDeliveryResponse struct {
	ID          string                               `json:"id"`
	Parcel      domain.Parcel                        `json:"parcel"`
	RiderId     string                               `json:"riderId"`
	CustomerId  string                               `json:"customerId"`
	Pickup      DetailedDeliveryResponseTimeAndPlace `json:"pickup"`
	Destination DetailedDeliveryResponseTimeAndPlace `json:"destination"`
	Status      int                                  `json:"status"`
	Route       domain.Line                          `json:"route"`
}

func CreateDetailedDeliveryResponse(delivery domain.Delivery) DetailedDeliveryResponse {
	return DetailedDeliveryResponse{
		ID:          delivery.ID,
		Parcel:      delivery.Parcel,
		RiderId:     delivery.RiderId,
		CustomerId:  delivery.CustomerId,
		Pickup:      DetailedDeliveryResponseTimeAndPlace(delivery.Pickup),
		Destination: DetailedDeliveryResponseTimeAndPlace(delivery.Destination),
		Status:      delivery.Status,
		Route:       delivery.Route,
	}
}
