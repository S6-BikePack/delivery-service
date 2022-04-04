package delivery_service

import (
	"delivery-service/internal/core/domain"
	"delivery-service/internal/core/ports"
	"errors"
	"github.com/google/uuid"
	"time"
)

type service struct {
	deliveryRepository ports.DeliveryRepository
	messagePublisher   ports.MessageBusPublisher
}

func New(deliveryRepository ports.DeliveryRepository, messagePublisher ports.MessageBusPublisher) *service {
	return &service{
		deliveryRepository: deliveryRepository,
		messagePublisher:   messagePublisher,
	}
}

func (srv *service) GetAll() ([]domain.Delivery, error) {
	return srv.deliveryRepository.GetAll()
}

func (srv *service) Get(id uuid.UUID) (domain.Delivery, error) {
	return srv.deliveryRepository.Get(id)
}

func (srv *service) Create(parcelId uuid.UUID, pickupPoint, deliveryPoint domain.Location, pickupTime time.Time) (domain.Delivery, error) {
	parcel, err := srv.deliveryRepository.GetParcel(parcelId)

	if err != nil {
		return domain.Delivery{}, errors.New("parcel not found with id:" + parcelId.String())
	}

	if parcel.Delivery != uuid.Nil {
		return domain.Delivery{}, errors.New("parcel is already part of a delivery")
	}

	delivery, err := domain.NewDelivery(parcel, pickupPoint, deliveryPoint, pickupTime)

	if err != nil {
		return domain.Delivery{}, err
	}

	delivery, err = srv.deliveryRepository.Save(delivery)

	if err != nil {
		return domain.Delivery{}, errors.New("saving new delivery failed")
	}

	srv.messagePublisher.CreateDelivery(delivery)
	return delivery, nil
}

func (srv *service) AssignRider(id, riderId uuid.UUID) (domain.Delivery, error) {
	delivery, err := srv.Get(id)

	if err != nil {
		return domain.Delivery{}, errors.New("could not find delivery with id: " + id.String())
	}

	if delivery.Rider != (domain.Rider{}) {
		return domain.Delivery{}, errors.New("delivery already has rider assigned")
	}

	rider, err := srv.deliveryRepository.GetRider(riderId)

	if err != nil {
		return domain.Delivery{}, errors.New("rider not found with id:" + riderId.String())
	}

	delivery.Rider = rider
	delivery.RiderId = riderId

	delivery, err = srv.deliveryRepository.Update(delivery)

	if err != nil {
		return domain.Delivery{}, err
	}

	srv.messagePublisher.UpdateDelivery(delivery)
	return delivery, nil
}

func (srv *service) StartDelivery(id uuid.UUID) (domain.Delivery, error) {
	delivery, err := srv.Get(id)

	if err != nil {
		return domain.Delivery{}, errors.New("could not find delivery with id")
	}

	delivery.Status = 2
	delivery.DeliveryTime = time.Now()

	delivery, err = srv.deliveryRepository.Update(delivery)

	if err != nil {
		return domain.Delivery{}, err
	}

	srv.messagePublisher.StartDelivery(delivery)
	return delivery, nil
}

func (srv *service) CompleteDelivery(id uuid.UUID) (domain.Delivery, error) {
	delivery, err := srv.Get(id)

	if err != nil {
		return domain.Delivery{}, errors.New("could not find delivery with id")
	}

	delivery.Status = 2
	delivery.DeliveryTime = time.Now()

	delivery, err = srv.deliveryRepository.Update(delivery)

	if err != nil {
		return domain.Delivery{}, err
	}

	srv.messagePublisher.CompleteDelivery(delivery)
	return delivery, nil
}

func (srv *service) GetRider(id uuid.UUID) (domain.Rider, error) {
	return srv.deliveryRepository.GetRider(id)
}

func (srv *service) SaveOrUpdateRider(rider domain.Rider) error {
	if rider.Name == "" || rider.ID == uuid.Nil {
		return errors.New("incomplete rider data")
	}

	_, err := srv.deliveryRepository.SaveOrUpdateRider(rider)

	return err
}

func (srv *service) SaveOrUpdateCustomer(customer domain.Customer) error {
	if customer.Name == "" || customer.LastName == "" || customer.ID == uuid.Nil {
		return errors.New("incomplete customer data")
	}

	_, err := srv.deliveryRepository.SaveOrUpdateCustomer(customer)

	return err
}
