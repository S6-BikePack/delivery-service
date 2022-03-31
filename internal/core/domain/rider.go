package domain

import (
	"github.com/google/uuid"
)

type Rider struct {
	ID   uuid.UUID `gorm:"type:uuid"`
	Name string
}

func NewRider(id uuid.UUID, name string) Rider {
	return Rider{ID: id, Name: name}
}
