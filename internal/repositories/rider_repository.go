package repositories

import (
	"delivery-service/internal/core/domain"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type riderRepository struct {
	Connection *gorm.DB
}

func NewRiderRepository(db *gorm.DB) (*riderRepository, error) {
	err := db.AutoMigrate(&domain.Rider{})

	if err != nil {
		return nil, err
	}

	database := riderRepository{
		Connection: db,
	}

	return &database, nil
}

func (repository *riderRepository) GetRider(riderId string) (domain.Rider, error) {
	var rider domain.Rider

	repository.Connection.Preload(clause.Associations).First(&rider, riderId)

	if (rider == domain.Rider{}) {
		return rider, errors.New("could not find rider")
	}

	return rider, nil
}

func (repository *riderRepository) CreateRider(rider domain.Rider) (domain.Rider, error) {
	if repository.Connection.Model(&rider).Where("id = ?", rider.ID).Updates(&rider).RowsAffected == 0 {
		repository.Connection.Create(&rider)
	}

	return repository.GetRider(rider.ID)
}

func (repository *riderRepository) UpdateRider(rider domain.Rider) (domain.Rider, error) {
	//TODO implement me
	panic("implement me")
}
