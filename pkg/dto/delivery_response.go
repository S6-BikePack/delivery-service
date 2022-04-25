package dto

import (
	"delivery-service/internal/core/domain"
	"time"
)

type DeliveryResponseParcel struct {
	ID     string            `json:"id"`
	Size   domain.Dimensions `json:"size"`
	Weight int               `json:"weight"`
}

type DeliveryResponseTimeAndPlace struct {
	Coordinates domain.Location `json:"coordinates"`
	Address     string          `json:"address"`
	Time        time.Time       `json:"time"`
}

type DeliveryResponse struct {
	ID          string                       `json:"id"`
	Parcel      DeliveryResponseParcel       `json:"parcel"`
	RiderId     string                       `json:"riderId"`
	CustomerId  string                       `json:"customerId"`
	Pickup      DeliveryResponseTimeAndPlace `json:"pickup"`
	Destination DeliveryResponseTimeAndPlace `json:"destination"`
	Status      int                          `json:"status"`
}

func CreateDeliveryResponse(delivery domain.Delivery) DeliveryResponse {
	return DeliveryResponse{
		ID: delivery.ID,
		Parcel: DeliveryResponseParcel{
			ID:     delivery.Parcel.ID,
			Size:   delivery.Parcel.Size,
			Weight: delivery.Parcel.Weight,
		},
		RiderId:    delivery.Rider.ID,
		CustomerId: delivery.Customer.ID,
		Pickup: DeliveryResponseTimeAndPlace{
			Coordinates: delivery.Pickup.Coordinates,
			Address:     delivery.Pickup.Address,
			Time:        delivery.Pickup.Time,
		},
		Destination: DeliveryResponseTimeAndPlace{
			Coordinates: delivery.Destination.Coordinates,
			Address:     delivery.Destination.Address,
			Time:        delivery.Destination.Time,
		},
		Status: delivery.Status,
	}
}
