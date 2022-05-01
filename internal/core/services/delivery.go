package services

import (
	"delivery-service/internal/core/domain"
	"delivery-service/internal/core/interfaces"
	"errors"
	"time"
)

type deliveryService struct {
	deliveryRepository interfaces.DeliveryRepository
	messagePublisher   interfaces.MessageBusPublisher
	riderService       interfaces.RiderService
}

func NewDeliveryService(deliveryRepository interfaces.DeliveryRepository, messagePublisher interfaces.MessageBusPublisher, riderService interfaces.RiderService) *deliveryService {
	return &deliveryService{
		deliveryRepository: deliveryRepository,
		messagePublisher:   messagePublisher,
		riderService:       riderService,
	}
}

func (srv *deliveryService) GetAll() ([]domain.Delivery, error) {
	return srv.deliveryRepository.GetAll()
}

func (srv *deliveryService) Get(id string) (domain.Delivery, error) {
	return srv.deliveryRepository.Get(id)
}

func (srv *deliveryService) GetByDistance(location domain.Location, radius int) []domain.Delivery {
	return srv.deliveryRepository.GetWithinRadius(location, radius)
}

func (srv *deliveryService) GetAroundRider(riderId string) ([]domain.Delivery, int) {
	rider, err := srv.riderService.Get(riderId)

	if err != nil || !rider.IsActive {
		return []domain.Delivery{}, 0
	}

	//TODO: Get radius based on amount of available riders in area

	radius := 1000

	deliveries := srv.deliveryRepository.GetWithinRadius(rider.Location, radius)

	for i, delivery := range deliveries {
		delivery.Pickup.Coordinates.Round()
		delivery.Destination.Coordinates.Round()

		deliveries[i] = delivery
	}

	return deliveries, radius
}

func (srv *deliveryService) Create(parcelId, ownerId string, pickup, destination domain.TimeAndPlace) (domain.Delivery, error) {
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

func (srv *deliveryService) AssignRider(id, riderId string) (domain.Delivery, error) {
	delivery, err := srv.Get(id)

	if err != nil {
		return domain.Delivery{}, errors.New("could not find delivery with id: " + id)
	}

	if delivery.Rider != (domain.Rider{}) {
		return domain.Delivery{}, errors.New("delivery already has rider assigned")
	}

	rider, err := srv.riderService.Get(riderId)

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

func (srv *deliveryService) StartDelivery(id string) (domain.Delivery, error) {
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

func (srv *deliveryService) CompleteDelivery(id string) (domain.Delivery, error) {
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

func (srv *deliveryService) SaveOrUpdateCustomer(customer domain.Customer) error {
	if customer.ID == "" {
		return errors.New("incomplete customer data")
	}

	_, err := srv.deliveryRepository.SaveOrUpdateCustomer(customer)

	return err
}

func (srv *deliveryService) SaveOrUpdateParcel(parcel domain.Parcel) error {
	if parcel.Size == (domain.Dimensions{}) || parcel.ID == "" {
		return errors.New("incomplete parcel data")
	}

	_, err := srv.deliveryRepository.SaveOrUpdateParcel(parcel)

	return err
}
