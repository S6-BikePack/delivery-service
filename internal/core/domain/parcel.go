package domain

import "github.com/google/uuid"

type Parcel struct {
	ID       uuid.UUID `gorm:"type:uuid"`
	Name     string
	Delivery uuid.UUID
}

func NewParcel(id uuid.UUID, name string) Parcel {
	return Parcel{ID: id, Name: name}
}
