package repositories

import (
	"delivery-service/internal/core/domain"
	"gorm.io/gorm"
)

type serviceAreaRepository struct {
	Connection *gorm.DB
}

func NewServiceAreaRepository(db *gorm.DB) (*serviceAreaRepository, error) {
	err := db.AutoMigrate(&domain.ServiceArea{})

	if err != nil {
		return nil, err
	}

	database := serviceAreaRepository{
		Connection: db,
	}

	db.Exec("INSERT INTO public.service_areas (id, identifier, rider_coverage) VALUES (0, 'Undefined', 0)")

	return &database, nil
}

func (repository *serviceAreaRepository) SaveOrUpdateServiceArea(serviceArea domain.ServiceArea) error {
	if repository.Connection.Model(&serviceArea).Where("id = ?", serviceArea.ID).Updates(&serviceArea).RowsAffected == 0 {
		create := repository.Connection.Create(&serviceArea)

		if create.Error != nil {
			return create.Error
		}
	}

	return nil
}
