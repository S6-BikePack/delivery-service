package domain

type Parcel struct {
	ID         string
	Name       string
	Size       Dimensions `gorm:"embedded"`
	Weight     int
	DeliveryId string
}

func NewParcel(id string, name string, size Dimensions, weight int) Parcel {
	return Parcel{ID: id, Name: name, Size: size, Weight: weight}
}
