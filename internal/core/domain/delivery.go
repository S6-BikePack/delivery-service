package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Delivery struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Parcel        Parcel    `gorm:"foreignKey:delivery"`
	RiderId       uuid.UUID
	Rider         Rider    `gorm:"foreignKey:RiderId"`
	PickupPoint   Location `gorm:"embedded;embeddedPrefix:pickup_"`
	PickupTime    time.Time
	DeliveryPoint Location `gorm:"embedded;embeddedPrefix:delivery_"`
	DeliveryTime  time.Time
	Status        int
}

func NewDelivery(parcel Parcel, pickupPoint Location, deliveryPoint Location, pickupTime time.Time) (Delivery, error) {

	if (parcel == Parcel{}) {
		return Delivery{}, errors.New("parcel can not be empty")
	}

	if (pickupTime == time.Time{}) {
		return Delivery{}, errors.New("pickup time can not be empty")
	}

	if (deliveryPoint == Location{}) {
		return Delivery{}, errors.New("location can not be empty")
	}

	delivery := Delivery{
		Parcel:        parcel,
		PickupPoint:   pickupPoint,
		PickupTime:    pickupTime,
		DeliveryPoint: deliveryPoint,
		Status:        0,
	}

	return delivery, nil
}
