package ports

import (
	"delivery-service/internal/core/domain"
)

type DeliveryRepository interface {
	GetAll() ([]domain.Delivery, error)
	Get(id string) (domain.Delivery, error)
	GetWithinRadius(location domain.Location, radius int) []domain.Delivery
	Save(delivery domain.Delivery) (domain.Delivery, error)
	Update(delivery domain.Delivery) (domain.Delivery, error)
	GetRider(riderId string) (domain.Rider, error)
	GetParcel(parcelId string) (domain.Parcel, error)
	GetCustomer(customerId string) (domain.Customer, error)
	SaveOrUpdateRider(rider domain.Rider) (domain.Rider, error)
	SaveOrUpdateCustomer(customer domain.Customer) (domain.Customer, error)
	SaveOrUpdateParcel(parcel domain.Parcel) (domain.Parcel, error)
}
