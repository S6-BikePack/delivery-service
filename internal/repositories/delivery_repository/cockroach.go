package delivery_repository

import (
	"delivery-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type cockroachdb struct {
	Connection *gorm.DB
}

func NewCockroachDB(connStr string) (*cockroachdb, error) {
	db, err := gorm.Open(postgres.Open(connStr))

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.Delivery{}, &domain.Rider{}, &domain.Parcel{})

	if err != nil {
		return nil, err
	}

	database := cockroachdb{
		Connection: db,
	}

	return &database, nil
}

func (repository *cockroachdb) Get(id uuid.UUID) (domain.Delivery, error) {
	var delivery domain.Delivery

	repository.Connection.Preload(clause.Associations).First(&delivery, id)

	return delivery, nil
}

func (repository *cockroachdb) GetAll() ([]domain.Delivery, error) {
	var deliveries []domain.Delivery

	repository.Connection.Find(&deliveries)

	return deliveries, nil
}

func (repository *cockroachdb) Save(delivery domain.Delivery) (domain.Delivery, error) {
	result := repository.Connection.Omit("RiderId").Create(&delivery)

	if result.Error != nil {
		return domain.Delivery{}, result.Error
	}

	return delivery, nil
}

func (repository *cockroachdb) Update(delivery domain.Delivery) (domain.Delivery, error) {
	result := repository.Connection.Model(&delivery).Updates(delivery)

	if result.Error != nil {
		return domain.Delivery{}, result.Error
	}

	return delivery, nil
}

func (repository *cockroachdb) GetRider(riderId uuid.UUID) domain.Rider {
	var rider domain.Rider

	repository.Connection.Preload(clause.Associations).First(&rider, riderId)

	return rider
}

func (repository *cockroachdb) SaveRider(rider domain.Rider) (domain.Rider, error) {
	result := repository.Connection.Create(&rider)

	if result.Error != nil {
		return domain.Rider{}, result.Error
	}

	return rider, nil
}

func (repository *cockroachdb) GetParcel(parcelId uuid.UUID) domain.Parcel {
	var parcel domain.Parcel

	repository.Connection.Preload(clause.Associations).First(&parcel, parcelId)

	return parcel
}

func (repository *cockroachdb) SaveParcel(parcel domain.Parcel) (domain.Parcel, error) {
	result := repository.Connection.Create(&parcel)

	if result.Error != nil {
		return domain.Parcel{}, result.Error
	}

	return parcel, nil
}
