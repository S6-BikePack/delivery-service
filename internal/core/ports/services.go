package ports

import (
	"delivery-service/internal/core/domain"
)

type DeliveryService interface {
	GetAll() ([]domain.Delivery, error)
	Get(id string) (domain.Delivery, error)
	GetByDistance(location domain.Location, radius int) []domain.Delivery
	Create(parcelId, ownerId string, pickup, destination domain.TimeAndPlace) (domain.Delivery, error)
	AssignRider(id, riderId string) (domain.Delivery, error)
	StartDelivery(id string) (domain.Delivery, error)
	CompleteDelivery(id string) (domain.Delivery, error)
	GetRider(id string) (domain.Rider, error)
	SaveOrUpdateRider(rider domain.Rider) error
	SaveOrUpdateCustomer(customer domain.Customer) error
	SaveOrUpdateParcel(parcel domain.Parcel) error
}
