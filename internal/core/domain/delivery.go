package domain

import (
	"errors"
	"time"
)

type Delivery struct {
	ID          string       `gorm:"default:uuid_generate_v4()" json:"id"`
	Parcel      Parcel       `gorm:"foreignKey:DeliveryId" json:"parcel"`
	RiderId     string       `json:"-"`
	Rider       Rider        `gorm:"foreignKey:RiderId" json:"rider"`
	CustomerId  string       `json:"-"`
	Customer    Customer     `gorm:"foreignKey:CustomerId" json:"customer"`
	Pickup      TimeAndPlace `gorm:"embedded;embeddedPrefix:pickup_" json:"pickup"`
	Destination TimeAndPlace `gorm:"embedded;embeddedPrefix:destination_" json:"destination"`
	Status      int          `json:"status"`
	Route       Line         `json:"route"`
}

func NewDelivery(parcel Parcel, owner Customer, pickup, destination TimeAndPlace) (Delivery, error) {
	if (parcel == Parcel{}) {
		return Delivery{}, errors.New("parcel can not be empty")
	}

	if (destination == TimeAndPlace{} || pickup == TimeAndPlace{}) {
		return Delivery{}, errors.New("pickup or destination locations can not be empty")
	}

	if pickup.Time.Before(time.Now()) {
		return Delivery{}, errors.New("pickup time can not be in the past")
	}

	delivery := Delivery{
		Parcel:      parcel,
		Customer:    owner,
		Pickup:      pickup,
		Destination: destination,
		Status:      0,
	}

	return delivery, nil
}
