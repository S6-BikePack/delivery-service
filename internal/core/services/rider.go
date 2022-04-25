package services

import (
	"delivery-service/internal/core/domain"
	"delivery-service/internal/core/interfaces"
	"errors"
)

type riderService struct {
	riderRepository  interfaces.RiderRepository
	messagePublisher interfaces.MessageBusPublisher
}

func NewRiderService(riderRepository interfaces.RiderRepository, messagePublisher interfaces.MessageBusPublisher) *riderService {
	return &riderService{
		riderRepository:  riderRepository,
		messagePublisher: messagePublisher,
	}
}

func (srv *riderService) Get(id string) (domain.Rider, error) {
	return srv.riderRepository.GetRider(id)
}

func (srv *riderService) Create(id, name string, serviceArea int) (domain.Rider, error) {
	if id == "" || name == "" {
		return domain.Rider{}, errors.New("incomplete rider data")
	}

	rider, err := srv.riderRepository.CreateRider(domain.NewRider(id, name, serviceArea))

	return rider, err
}

func (srv *riderService) Update(rider domain.Rider) (domain.Rider, error) {
	existing, err := srv.Get(rider.ID)

	if err != nil {
		return domain.Rider{}, err
	}

	if rider.Name != "" {
		existing.Name = rider.Name
	}

	if rider.ServiceArea != 0 {
		existing.ServiceArea = rider.ServiceArea
	}

	return srv.riderRepository.UpdateRider(existing)
}

func (srv *riderService) UpdateActiveStatus(id string, status bool) (domain.Rider, error) {
	existing, err := srv.Get(id)

	if err != nil {
		return domain.Rider{}, err
	}

	if existing.IsActive == status {
		return existing, nil
	}

	existing.IsActive = status

	return srv.riderRepository.UpdateRider(existing)
}

func (srv *riderService) UpdateLocation(id string, location domain.Location) error {
	existing, err := srv.Get(id)

	if err != nil {
		return err
	}

	existing.Location = location

	_, err = srv.riderRepository.UpdateRider(existing)

	return err
}
