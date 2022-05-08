package mocks

import (
	"delivery-service/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type DeliveryRepository struct {
	mock.Mock
}

func (d *DeliveryRepository) GetAll() ([]domain.Delivery, error) {
	args := d.Called()
	return args.Get(0).([]domain.Delivery), args.Error(1)
}

func (d *DeliveryRepository) Get(id string) (domain.Delivery, error) {
	args := d.Called(id)
	return args.Get(0).(domain.Delivery), args.Error(1)
}

func (d *DeliveryRepository) GetWithinRadius(location domain.Location, radius int) []domain.Delivery {
	args := d.Called(location, radius)
	return args.Get(0).([]domain.Delivery)
}

func (d *DeliveryRepository) Save(delivery domain.Delivery) (domain.Delivery, error) {
	args := d.Called(delivery)
	return args.Get(0).(domain.Delivery), args.Error(1)
}

func (d *DeliveryRepository) Update(delivery domain.Delivery) (domain.Delivery, error) {
	args := d.Called(delivery)
	return args.Get(0).(domain.Delivery), args.Error(1)
}

func (d *DeliveryRepository) GetRider(riderId string) (domain.Rider, error) {
	args := d.Called(riderId)
	return args.Get(0).(domain.Rider), args.Error(1)
}

func (d *DeliveryRepository) GetParcel(parcelId string) (domain.Parcel, error) {
	args := d.Called(parcelId)
	return args.Get(0).(domain.Parcel), args.Error(1)
}

func (d *DeliveryRepository) GetCustomer(customerId string) (domain.Customer, error) {
	args := d.Called(customerId)
	return args.Get(0).(domain.Customer), args.Error(1)
}

func (d *DeliveryRepository) CreateRider(rider domain.Rider) (domain.Rider, error) {
	args := d.Called(rider)
	return args.Get(0).(domain.Rider), args.Error(1)
}

func (d *DeliveryRepository) SaveOrUpdateCustomer(customer domain.Customer) (domain.Customer, error) {
	args := d.Called(customer)
	return args.Get(0).(domain.Customer), args.Error(1)
}

func (d *DeliveryRepository) SaveOrUpdateParcel(parcel domain.Parcel) (domain.Parcel, error) {
	args := d.Called(parcel)
	return args.Get(0).(domain.Parcel), args.Error(1)
}

func (d *DeliveryRepository) UpdateRider(rider domain.Rider) (domain.Rider, error) {
	args := d.Called(rider)
	return args.Get(0).(domain.Rider), args.Error(1)
}
