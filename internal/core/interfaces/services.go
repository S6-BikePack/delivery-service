package interfaces

import (
	"delivery-service/internal/core/domain"
)

type DeliveryService interface {
	GetAll() ([]domain.Delivery, error)
	Get(id string) (domain.Delivery, error)
	GetByDistance(location domain.Location, radius int) []domain.Delivery
	GetAroundRider(riderId string) ([]domain.Delivery, int)
	Create(parcelId, ownerId string, pickup, destination domain.TimeAndPlace) (domain.Delivery, error)
	AssignRider(id, riderId string) (domain.Delivery, error)
	StartDelivery(id string) (domain.Delivery, error)
	CompleteDelivery(id string) (domain.Delivery, error)
	SaveOrUpdateCustomer(customer domain.Customer) error
	SaveOrUpdateParcel(parcel domain.Parcel) error
}

type RiderService interface {
	Get(id string) (domain.Rider, error)
	Create(id, name string, serviceArea int) (domain.Rider, error)
	Update(rider domain.Rider) (domain.Rider, error)
	UpdateActiveStatus(id string, status bool) (domain.Rider, error)
	UpdateLocation(id string, location domain.Location) error
}

type ServiceAreaService interface {
	SaveOrUpdateServiceArea(serviceArea domain.ServiceArea) error
}
