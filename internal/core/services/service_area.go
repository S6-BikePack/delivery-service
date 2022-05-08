package services

import (
	"delivery-service/internal/core/domain"
	"delivery-service/internal/core/interfaces"
)

type serviceAreaService struct {
	serviceAreaRepository interfaces.ServiceAreaRepository
}

func NewServiceAreaService(serviceAreaRepository interfaces.ServiceAreaRepository) *serviceAreaService {
	return &serviceAreaService{
		serviceAreaRepository: serviceAreaRepository,
	}
}

func (s *serviceAreaService) SaveOrUpdateServiceArea(serviceArea domain.ServiceArea) error {
	return s.serviceAreaRepository.SaveOrUpdateServiceArea(serviceArea)
}
