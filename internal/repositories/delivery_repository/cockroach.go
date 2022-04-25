package delivery_repository

import (
	"delivery-service/internal/core/domain"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type cockroachdb struct {
	Connection *gorm.DB
}

func NewCockroachDB(db *gorm.DB) (*cockroachdb, error) {
	err := db.AutoMigrate(&domain.Delivery{}, &domain.Rider{}, &domain.Parcel{}, &domain.Customer{})

	if err != nil {
		return nil, err
	}

	database := cockroachdb{
		Connection: db,
	}

	return &database, nil
}

func (repository *cockroachdb) Get(id string) (domain.Delivery, error) {
	var delivery domain.Delivery

	repository.Connection.Preload(clause.Associations).First(&delivery, "id = ?", id)

	return delivery, nil
}

func (repository *cockroachdb) GetAll() ([]domain.Delivery, error) {
	var deliveries []domain.Delivery

	repository.Connection.Find(&deliveries)

	return deliveries, nil
}

func (repository *cockroachdb) GetWithinRadius(location domain.Location, radius int) []domain.Delivery {
	var deliveries []domain.Delivery

	repository.Connection.Preload(clause.Associations).Where("st_distancesphere(pickup_coordinates, ST_MakePoint(?, ?)) <= ?", location.Longitude, location.Latitude, radius).Find(&deliveries)

	if deliveries == nil {
		return []domain.Delivery{}
	}

	return deliveries
}

func (repository *cockroachdb) Save(delivery domain.Delivery) (domain.Delivery, error) {

	result := repository.Connection.Omit("RiderId", "Rider").Create(&delivery)

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

func (repository *cockroachdb) GetRider(riderId string) (domain.Rider, error) {
	var rider domain.Rider

	repository.Connection.Preload(clause.Associations).First(&rider, riderId)

	if (rider == domain.Rider{}) {
		return rider, errors.New("could not find rider")
	}

	return rider, nil
}

func (repository *cockroachdb) SaveOrUpdateRider(rider domain.Rider) (domain.Rider, error) {
	if repository.Connection.Model(&rider).Where("id = ?", rider.ID).Updates(&rider).RowsAffected == 0 {
		repository.Connection.Create(&rider)
	}

	return repository.GetRider(rider.ID)
}

func (repository *cockroachdb) GetParcel(parcelId string) (domain.Parcel, error) {
	var parcel domain.Parcel

	repository.Connection.First(&parcel, "id = ?", parcelId)

	if (parcel == domain.Parcel{}) {
		return parcel, errors.New("could not find customer")
	}

	return parcel, nil
}

func (repository *cockroachdb) SaveParcel(parcel domain.Parcel) (domain.Parcel, error) {
	result := repository.Connection.Create(&parcel)

	if result.Error != nil {
		return domain.Parcel{}, result.Error
	}

	return parcel, nil
}

func (repository *cockroachdb) GetCustomer(customerId string) (domain.Customer, error) {
	var customer domain.Customer

	repository.Connection.Preload(clause.Associations).First(&customer, "id = ?", customerId)

	if (customer == domain.Customer{}) {
		return customer, errors.New("could not find customer")
	}

	return customer, nil
}

func (repository *cockroachdb) SaveOrUpdateCustomer(customer domain.Customer) (domain.Customer, error) {
	if repository.Connection.Model(&customer).Where("id = ?", customer.ID).Updates(&customer).RowsAffected == 0 {
		repository.Connection.Create(&customer)
	}

	return repository.GetCustomer(customer.ID)
}

func (repository *cockroachdb) SaveOrUpdateParcel(parcel domain.Parcel) (domain.Parcel, error) {
	if repository.Connection.Model(&parcel).Omit("delivery_id").Where("id = ?", parcel.ID).Updates(&parcel).RowsAffected == 0 {
		repository.Connection.Omit("delivery_id").Create(&parcel)
	}

	return repository.GetParcel(parcel.ID)
}
