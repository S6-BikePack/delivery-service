package domain

type ServiceArea struct {
	ID            int
	Identifier    string
	RiderCoverage int
}

func NewServiceArea(id int, identifier string) ServiceArea {
	return ServiceArea{
		ID:         id,
		Identifier: identifier,
	}
}
