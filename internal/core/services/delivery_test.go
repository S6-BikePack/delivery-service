package services

import (
	"delivery-service/internal/core/domain"
	"delivery-service/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func MockDeliveries() []domain.Delivery {
	return []domain.Delivery{
		{
			ID: "delivery-1",
			Parcel: domain.Parcel{
				ID: "parcel-1",
				Size: domain.Dimensions{
					Width:  10,
					Height: 10,
					Depth:  10,
				},
				Weight:      10,
				ServiceArea: 1,
			},
			Rider: domain.Rider{
				ID:          "rider-1",
				ServiceArea: 1,
			},
			Customer: domain.Customer{
				ID:          "customer-1",
				ServiceArea: 1,
			},
			Pickup: domain.TimeAndPlace{
				Coordinates: domain.Location{
					Latitude:  1.123451242,
					Longitude: 2.123415342,
				},
				Address: "Test street 1",
				Time:    time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
			},
			Destination: domain.TimeAndPlace{
				Coordinates: domain.Location{
					Latitude:  2.534543634,
					Longitude: 3.346123423,
				},
				Address: "Test street 2",
			},
			Status: 0,
		},
		{
			ID: "delivery-2",
			Parcel: domain.Parcel{
				ID: "parcel-2",
				Size: domain.Dimensions{
					Width:  20,
					Height: 20,
					Depth:  20,
				},
				Weight:      20,
				ServiceArea: 2,
			},
			Rider: domain.Rider{
				ID:          "rider-2",
				ServiceArea: 2,
			},
			Customer: domain.Customer{
				ID:          "customer-1",
				ServiceArea: 2,
			},
			Pickup: domain.TimeAndPlace{
				Coordinates: domain.Location{
					Latitude:  1.67354234,
					Longitude: 2.23523611,
				},
				Address: "Test street 1",
				Time:    time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
			},
			Destination: domain.TimeAndPlace{
				Coordinates: domain.Location{
					Latitude:  2.64533143,
					Longitude: 3.53141233,
				},
				Address: "Test street 2",
			},
			Status: 1,
		},
	}
}

type DeliveryServiceTestSuite struct {
	suite.Suite
}

func (suite *DeliveryServiceTestSuite) Test_GetAll() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	riderService := NewRiderService(mockRepository, mockMessageBus)
	sut := NewDeliveryService(mockRepository, mockMessageBus, riderService)

	mockRepository.On("GetAll").Return(MockDeliveries(), nil)

	deliveries, err := sut.GetAll()

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), MockDeliveries(), deliveries)
}

func (suite *DeliveryServiceTestSuite) Test_Get() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	riderService := NewRiderService(mockRepository, mockMessageBus)
	sut := NewDeliveryService(mockRepository, mockMessageBus, riderService)
	id := "delivery-1"

	mockRepository.On("Get", id).Return(domain.Delivery{ID: id}, nil)

	sut.Get(id)

	mockRepository.AssertExpectations(suite.T())
}

func (suite *DeliveryServiceTestSuite) Test_GetByDistance() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	riderService := NewRiderService(mockRepository, mockMessageBus)
	sut := NewDeliveryService(mockRepository, mockMessageBus, riderService)

	location := domain.Location{
		Latitude:  1,
		Longitude: 2,
	}
	radius := 100

	mockRepository.On("GetWithinRadius", location, radius).Return([]domain.Delivery{}, nil)

	sut.GetByDistance(location, radius)

	mockRepository.AssertExpectations(suite.T())
}

func (suite *DeliveryServiceTestSuite) Test_GetAroundRider() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	riderService := NewRiderService(mockRepository, mockMessageBus)
	sut := NewDeliveryService(mockRepository, mockMessageBus, riderService)

	rider := domain.Rider{
		ID:          "test-rider",
		Name:        "Test",
		ServiceArea: 1,
		IsActive:    true,
		Location: domain.Location{
			Latitude:  1,
			Longitude: 2,
		},
	}

	roundedCoordinates := domain.Location{
		Latitude:  1.124,
		Longitude: 2.124,
	}

	mockRepository.On("GetRider", rider.ID).Return(rider, nil)
	mockRepository.On("GetWithinRadius", rider.Location, 1000).Return(MockDeliveries(), nil)

	deliveries, radius := sut.GetAroundRider(rider.ID)

	mockRepository.AssertExpectations(suite.T())

	assert.Equal(suite.T(), 1000, radius)
	assert.Equal(suite.T(), len(MockDeliveries()), len(deliveries))
	assert.Equal(suite.T(), roundedCoordinates, deliveries[0].Pickup.Coordinates)
}

func (suite *DeliveryServiceTestSuite) Test_GetAroundRider_NotFound() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	riderService := NewRiderService(mockRepository, mockMessageBus)
	sut := NewDeliveryService(mockRepository, mockMessageBus, riderService)

	rider := domain.Rider{
		ID: "test-rider",
	}

	mockRepository.On("GetRider", rider.ID).Return(domain.Rider{}, errors.New("rider not found"))

	deliveries, radius := sut.GetAroundRider(rider.ID)

	mockRepository.AssertExpectations(suite.T())

	assert.Equal(suite.T(), []domain.Delivery{}, deliveries)
	assert.Equal(suite.T(), 0, radius)
}

func TestDeliveryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DeliveryServiceTestSuite))
}
