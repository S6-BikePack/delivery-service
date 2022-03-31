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
	GetRider(riderId uuid.UUID) domain.Rider
	SaveRider(rider domain.Rider) (domain.Rider, error)
	GetParcel(parcelId uuid.UUID) domain.Parcel
	SaveParcel(parcel domain.Parcel) (domain.Parcel, error)
}
