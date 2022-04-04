package domain

import "github.com/google/uuid"

type Customer struct {
	ID       uuid.UUID `gorm:"type:uuid"`
	Name     string
	LastName string
}

func NewCustomer(id uuid.UUID, name, lastName string) Customer {
	return Customer{ID: id, Name: name, LastName: lastName}
}
