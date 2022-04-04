package ports

import (
	"delivery-service/internal/core/domain"
	"github.com/google/uuid"
	"time"
)

type DeliveryService interface {
	GetAll() ([]domain.Delivery, error)
	Get(id uuid.UUID) (domain.Delivery, error)
	Create(parcelId uuid.UUID, pickupPoint, deliveryPoint domain.Location, pickupTime time.Time) (domain.Delivery, error)
	AssignRider(id, riderId uuid.UUID) (domain.Delivery, error)
	StartDelivery(id uuid.UUID) (domain.Delivery, error)
	CompleteDelivery(id uuid.UUID) (domain.Delivery, error)
	GetRider(id uuid.UUID) (domain.Rider, error)
	SaveOrUpdateRider(rider domain.Rider) error
	SaveOrUpdateCustomer(customer domain.Customer) error
}
