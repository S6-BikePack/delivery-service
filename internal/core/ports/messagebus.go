package ports

import (
	"delivery-service/internal/core/domain"
)

type MessageBusPublisher interface {
	CreateDelivery(delivery domain.Delivery) error
	UpdateDelivery(delivery domain.Delivery) error
}
