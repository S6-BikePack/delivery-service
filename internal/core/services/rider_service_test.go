package services

import (
	"delivery-service/internal/core/domain"
	"delivery-service/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RiderServiceTestSuite struct {
	suite.Suite
}

func (suite *RiderServiceTestSuite) Test_Get() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
	}

	mockRepository.On("GetRider", rider.ID).Return(rider, nil)

	result, err := sut.Get(rider.ID)

	mockRepository.AssertExpectations(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), rider, result)
}

func (suite *RiderServiceTestSuite) Test_Create() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID:   "rider-1",
		Name: "test",
	}

	mockRepository.On("CreateRider", rider).Return(rider, nil)

	result, err := sut.Create(rider.ID, rider.Name, 0)

	mockRepository.AssertExpectations(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), rider, result)
}

func (suite *RiderServiceTestSuite) Test_Create_MissingId() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
	}

	_, err := sut.Create("", rider.Name, 1)

	assert.Error(suite.T(), err)
}

func (suite *RiderServiceTestSuite) Test_Create_MissingName() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
	}

	_, err := sut.Create(rider.ID, "", 1)

	assert.Error(suite.T(), err)
}

func (suite *RiderServiceTestSuite) Test_Update() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	oldRider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
	}

	rider := domain.Rider{
		ID:          "rider-1",
		Name:        "newTest",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
	}

	mockRepository.On("GetRider", rider.ID).Return(oldRider, nil)
	mockRepository.On("UpdateRider", rider).Return(rider, nil)

	result, err := sut.Update(rider)

	mockRepository.AssertExpectations(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), rider, result)
}

func (suite *RiderServiceTestSuite) Test_Update_NameOnly() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	oldRider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
	}

	rider := domain.Rider{
		ID:   "rider-1",
		Name: "newTest",
	}

	combined := domain.Rider{
		ID:          oldRider.ID,
		Name:        rider.Name,
		ServiceArea: oldRider.ServiceArea,
	}

	mockRepository.On("GetRider", rider.ID).Return(oldRider, nil)
	mockRepository.On("UpdateRider", combined).Return(combined, nil)

	result, err := sut.Update(rider)

	mockRepository.AssertExpectations(suite.T())

	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), rider.Name, result.Name)
	assert.Equal(suite.T(), oldRider.ServiceArea, result.ServiceArea)
}

func (suite *RiderServiceTestSuite) Test_Update_ServiceAreaOnly() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	oldRider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
	}

	rider := domain.Rider{
		ID:          "rider-1",
		ServiceArea: domain.ServiceArea{ID: 2, Identifier: "test-2"},
	}

	combined := domain.Rider{
		ID:          oldRider.ID,
		Name:        oldRider.Name,
		ServiceArea: rider.ServiceArea,
	}

	mockRepository.On("GetRider", rider.ID).Return(oldRider, nil)
	mockRepository.On("UpdateRider", combined).Return(combined, nil)

	result, err := sut.Update(rider)

	mockRepository.AssertExpectations(suite.T())

	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), oldRider.Name, result.Name)
	assert.Equal(suite.T(), rider.ServiceArea, result.ServiceArea)
}

func (suite *RiderServiceTestSuite) Test_Update_NotFound() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
	}

	mockRepository.On("GetRider", rider.ID).Return(domain.Rider{}, errors.New("rider not found"))

	_, err := sut.Update(rider)

	assert.Error(suite.T(), err)
}

func (suite *RiderServiceTestSuite) Test_UpdateActiveStatus() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
		IsActive:    false,
	}

	updatedRider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
		IsActive:    true,
	}

	mockRepository.On("GetRider", rider.ID).Return(rider, nil)
	mockRepository.On("UpdateRider", updatedRider).Return(updatedRider, nil)

	result, err := sut.UpdateActiveStatus(rider.ID, true)

	mockRepository.AssertExpectations(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedRider, result)
}

func (suite *RiderServiceTestSuite) Test_UpdateActiveStatus_NotFound() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID: "rider-1",
	}

	mockRepository.On("GetRider", rider.ID).Return(domain.Rider{}, errors.New("rider not found"))

	_, err := sut.UpdateActiveStatus(rider.ID, true)

	mockRepository.AssertExpectations(suite.T())

	assert.Error(suite.T(), err)
}

func (suite *RiderServiceTestSuite) Test_UpdateActiveStatus_NotUpdated() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
		IsActive:    false,
	}

	updatedRider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
		IsActive:    true,
	}

	mockRepository.On("GetRider", rider.ID).Return(rider, nil)
	mockRepository.On("UpdateRider", updatedRider).Return(domain.Rider{}, errors.New("could not update"))

	_, err := sut.UpdateActiveStatus(rider.ID, true)

	mockRepository.AssertExpectations(suite.T())

	assert.Error(suite.T(), err)
}

func (suite *RiderServiceTestSuite) Test_UpdateActiveStatus_Same() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID:          "rider-1",
		Name:        "test",
		ServiceArea: domain.ServiceArea{ID: 1, Identifier: "test"},
		IsActive:    true,
	}

	mockRepository.On("GetRider", rider.ID).Return(rider, nil)

	result, err := sut.UpdateActiveStatus(rider.ID, true)

	mockRepository.AssertExpectations(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), rider, result)
}

func (suite *RiderServiceTestSuite) Test_UpdateLocation() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID: "rider-1",
	}

	updatedRider := domain.Rider{
		ID: "rider-1",
		Location: domain.Location{
			Latitude:  1,
			Longitude: 1,
		},
	}

	mockRepository.On("GetRider", rider.ID).Return(rider, nil)
	mockRepository.On("UpdateRider", updatedRider).Return(updatedRider, nil)

	err := sut.UpdateLocation(rider.ID, updatedRider.Location)

	mockRepository.AssertExpectations(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *RiderServiceTestSuite) Test_UpdateLocation_NotFound() {
	mockRepository := new(mocks.DeliveryRepository)
	mockMessageBus := new(mocks.MessageBusPublisher)

	sut := NewRiderService(mockRepository, mockMessageBus)

	rider := domain.Rider{
		ID: "rider-1",
		Location: domain.Location{
			Latitude:  1,
			Longitude: 1,
		},
	}

	mockRepository.On("GetRider", rider.ID).Return(domain.Rider{}, errors.New("rider not found"))

	err := sut.UpdateLocation(rider.ID, rider.Location)

	mockRepository.AssertExpectations(suite.T())

	assert.Error(suite.T(), err)
}

func TestRiderServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RiderServiceTestSuite))
}
