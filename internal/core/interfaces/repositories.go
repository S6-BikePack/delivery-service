package interfaces

import (
	"delivery-service/internal/core/domain"
)

type DeliveryRepository interface {
	GetAll() ([]domain.Delivery, error)
	Get(id string) (domain.Delivery, error)
	GetWithinRadius(location domain.Location, radius int) []domain.Delivery
	Save(delivery domain.Delivery) (domain.Delivery, error)
	Update(delivery domain.Delivery) (domain.Delivery, error)

	GetParcel(parcelId string) (domain.Parcel, error)
	SaveOrUpdateParcel(parcel domain.Parcel) (domain.Parcel, error)

	GetCustomer(customerId string) (domain.Customer, error)
	SaveOrUpdateCustomer(customer domain.Customer) (domain.Customer, error)

	//GetNearRider(riderId string) ([]domain.Delivery, error)
}

type RiderRepository interface {
	GetRider(riderId string) (domain.Rider, error)
	CreateRider(rider domain.Rider) (domain.Rider, error)
	UpdateRider(rider domain.Rider) (domain.Rider, error)
}
