package mocks

import (
	"delivery-service/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MessageBusPublisher struct {
	mock.Mock
}

func (m *MessageBusPublisher) CreateDelivery(delivery domain.Delivery) error {
	m.Called(delivery)
	return nil
}

func (m *MessageBusPublisher) UpdateDelivery(delivery domain.Delivery) error {
	m.Called(delivery)
	return nil
}

func (m *MessageBusPublisher) StartDelivery(delivery domain.Delivery) error {
	m.Called(delivery)
	return nil
}

func (m *MessageBusPublisher) CompleteDelivery(delivery domain.Delivery) error {
	m.Called(delivery)
	return nil
}
