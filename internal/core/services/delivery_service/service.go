package delivery_service

import (
	"delivery-service/internal/core/domain"
	"delivery-service/internal/core/ports"
	"errors"
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

func (srv *service) Get(id string) (domain.Delivery, error) {
	return srv.deliveryRepository.Get(id)
}

func (srv *service) GetByDistance(location domain.Location, radius int) []domain.Delivery {
	return srv.deliveryRepository.GetWithinRadius(location, radius)
}

func (srv *service) Create(parcelId, ownerId string, pickup, destination domain.TimeAndPlace) (domain.Delivery, error) {
	owner, err := srv.deliveryRepository.GetCustomer(ownerId)

	if err != nil {
		return domain.Delivery{}, err
	}

	parcel, err := srv.deliveryRepository.GetParcel(parcelId)

	if err != nil {
		return domain.Delivery{}, errors.New("parcel not found with id:" + parcelId)
	}

	if parcel.DeliveryId != "" {
		return domain.Delivery{}, errors.New("parcel is already part of a delivery")
	}

	delivery, err := domain.NewDelivery(parcel, owner, pickup, destination)

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

func (srv *service) AssignRider(id, riderId string) (domain.Delivery, error) {
	delivery, err := srv.Get(id)

	if err != nil {
		return domain.Delivery{}, errors.New("could not find delivery with id: " + id)
	}

	if delivery.Rider != (domain.Rider{}) {
		return domain.Delivery{}, errors.New("delivery already has rider assigned")
	}

	rider, err := srv.deliveryRepository.GetRider(riderId)

	if err != nil {
		return domain.Delivery{}, errors.New("rider not found with id:" + riderId)
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

func (srv *service) StartDelivery(id string) (domain.Delivery, error) {
	delivery, err := srv.Get(id)

	if err != nil {
		return domain.Delivery{}, errors.New("could not find delivery with id")
	}

	delivery.Status = 2

	delivery, err = srv.deliveryRepository.Update(delivery)

	if err != nil {
		return domain.Delivery{}, err
	}

	srv.messagePublisher.StartDelivery(delivery)
	return delivery, nil
}

func (srv *service) CompleteDelivery(id string) (domain.Delivery, error) {
	delivery, err := srv.Get(id)

	if err != nil {
		return domain.Delivery{}, errors.New("could not find delivery with id")
	}

	delivery.Status = 3
	delivery.Destination.Time = time.Now()

	delivery, err = srv.deliveryRepository.Update(delivery)

	if err != nil {
		return domain.Delivery{}, err
	}

	srv.messagePublisher.CompleteDelivery(delivery)
	return delivery, nil
}

func (srv *service) GetRider(id string) (domain.Rider, error) {
	return srv.deliveryRepository.GetRider(id)
}

func (srv *service) SaveOrUpdateRider(rider domain.Rider) error {
	if rider.ID == "" {
		return errors.New("incomplete rider data")
	}

	_, err := srv.deliveryRepository.SaveOrUpdateRider(rider)

	return err
}

func (srv *service) SaveOrUpdateCustomer(customer domain.Customer) error {
	if customer.ID == "" {
		return errors.New("incomplete customer data")
	}

	_, err := srv.deliveryRepository.SaveOrUpdateCustomer(customer)

	return err
}

func (srv *service) SaveOrUpdateParcel(parcel domain.Parcel) error {
	if parcel.Size == (domain.Dimensions{}) || parcel.ID == "" {
		return errors.New("incomplete parcel data")
	}

	_, err := srv.deliveryRepository.SaveOrUpdateParcel(parcel)

	return err
}
