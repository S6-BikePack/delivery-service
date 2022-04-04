package ports

import (
	"delivery-service/internal/core/domain"
	"github.com/google/uuid"
)

type DeliveryRepository interface {
	GetAll() ([]domain.Delivery, error)
	Get(id uuid.UUID) (domain.Delivery, error)
	Save(delivery domain.Delivery) (domain.Delivery, error)
	Update(delivery domain.Delivery) (domain.Delivery, error)
	GetRider(riderId uuid.UUID) (domain.Rider, error)
	SaveOrUpdateRider(rider domain.Rider) (domain.Rider, error)
	GetParcel(parcelId uuid.UUID) (domain.Parcel, error)
	SaveParcel(parcel domain.Parcel) (domain.Parcel, error)
	GetCustomer(customerId uuid.UUID) (domain.Customer, error)
	SaveOrUpdateCustomer(customer domain.Customer) (domain.Customer, error)
}
