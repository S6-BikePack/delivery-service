package domain

import (
	"errors"
	"time"
)

type Delivery struct {
	ID            string
	ParcelId      string
	Parcel        Parcel `gorm:"foreignKey:DeliveryId"`
	RiderId       string
	Rider         Rider `gorm:"foreignKey:RiderId"`
	CustomerId    string
	Customer      Customer `gorm:"foreignKey:CustomerId"`
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
