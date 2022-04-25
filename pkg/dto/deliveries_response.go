package dto

import (
	"delivery-service/internal/core/domain"
	"time"
)

type DeliveriesResponseParcel struct {
	Size   domain.Dimensions `json:"size"`
	Weight int               `json:"weight"`
}

type DeliveriesResponseTimeAndPlace struct {
	Coordinates domain.Location `json:"coordinates"`
	Time        time.Time       `json:"time"`
}

type DeliveriesResponse struct {
	ID          string                         `json:"id"`
	Parcel      DeliveriesResponseParcel       `json:"parcel"`
	Pickup      DeliveriesResponseTimeAndPlace `json:"pickup"`
	Destination DeliveriesResponseTimeAndPlace `json:"destination"`
	Status      int                            `json:"status"`
}

func CreateDeliveriesResponse(delivery domain.Delivery) DeliveriesResponse {
	return DeliveriesResponse{
		ID: delivery.ID,
		Parcel: DeliveriesResponseParcel{
			Size:   delivery.Parcel.Size,
			Weight: delivery.Parcel.Weight,
		},
		Pickup: DeliveriesResponseTimeAndPlace{
			Coordinates: delivery.Pickup.Coordinates,
			Time:        delivery.Pickup.Time,
		},
		Destination: DeliveriesResponseTimeAndPlace{
			Coordinates: delivery.Destination.Coordinates,
			Time:        delivery.Destination.Time,
		},
		Status: delivery.Status,
	}
}

type DeliveryListResponse []*DeliveriesResponse

func CreateDeliveryListResponse(deliveries []domain.Delivery) DeliveryListResponse {
	deliveriesResponse := DeliveryListResponse{}
	for _, d := range deliveries {
		delivery := CreateDeliveriesResponse(d)
		deliveriesResponse = append(deliveriesResponse, &delivery)
	}
	return deliveriesResponse
}
