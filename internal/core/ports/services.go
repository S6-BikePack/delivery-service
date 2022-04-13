package ports

import (
	"delivery-service/internal/core/domain"
	"time"
)

type DeliveryService interface {
	GetAll() ([]domain.Delivery, error)
	Get(id string) (domain.Delivery, error)
	Create(parcelId string, pickupPoint, deliveryPoint domain.Location, pickupTime time.Time) (domain.Delivery, error)
	AssignRider(id, riderId string) (domain.Delivery, error)
	StartDelivery(id string) (domain.Delivery, error)
	CompleteDelivery(id string) (domain.Delivery, error)
	GetRider(id string) (domain.Rider, error)
	SaveOrUpdateRider(rider domain.Rider) error
	SaveOrUpdateCustomer(customer domain.Customer) error
	SaveOrUpdateParcel(parcel domain.Parcel) error
}
