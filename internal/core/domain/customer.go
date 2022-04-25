package domain

type Customer struct {
	ID          string `json:"id"`
	ServiceArea int    `json:"-"`
}

func NewCustomer(id string, serviceArea int) Customer {
	return Customer{ID: id, ServiceArea: serviceArea}
}
