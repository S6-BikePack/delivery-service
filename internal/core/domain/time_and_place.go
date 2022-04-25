package domain

import (
	"errors"
	"time"
)

type TimeAndPlace struct {
	Coordinates Location  `gorm:"embedded" json:"coordinates"`
	Address     string    `json:"address"`
	Time        time.Time `json:"time"`
}

func NewTimeAndPlace(address string, coordinates Location) (TimeAndPlace, error) {

	if address == "" {
		return TimeAndPlace{}, errors.New("address cannot be empty")
	}

	if coordinates == (Location{}) {
		return TimeAndPlace{}, errors.New("coordinates cannot be empty")
	}

	return TimeAndPlace{
		Coordinates: coordinates,
		Address:     address,
	}, nil
}
