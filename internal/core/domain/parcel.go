package domain

type Parcel struct {
	ID        string
	Name      string
	Size      Dimensions
	DeliverId string
}

func NewParcel(id string, name string, size Dimensions) Parcel {
	return Parcel{ID: id, Name: name, Size: size}
}
