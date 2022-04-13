package domain

type Customer struct {
	ID          string
	Name        string
	ServiceArea int
}

func NewCustomer(id string, name string, serviceArea int) Customer {
	return Customer{ID: id, Name: name, ServiceArea: serviceArea}
}
