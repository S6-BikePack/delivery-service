package domain

type Rider struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	ServiceArea int      `json:"serviceArea"`
	IsActive    bool     `json:"isActive"`
	Location    Location `gorm:"embedded" json:"location"`
}

func NewRider(id, name string, serviceArea int) Rider {
	return Rider{ID: id, Name: name, ServiceArea: serviceArea, IsActive: false}
}
