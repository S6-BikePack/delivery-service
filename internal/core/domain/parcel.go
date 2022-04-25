package domain

type Parcel struct {
	ID          string     `json:"id"`
	Size        Dimensions `gorm:"embedded" json:"size"`
	Weight      int        `json:"weight"`
	ServiceArea int        `json:"serviceArea"`
	DeliveryId  string     `json:"-"`
}

func NewParcel(id string, size Dimensions, weight int) Parcel {
	return Parcel{ID: id, Size: size, Weight: weight}
}
