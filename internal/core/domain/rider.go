package domain

type Rider struct {
	ID          string `json:"id"`
	ServiceArea int    `json:"serviceArea"`
}

func NewRider(id string) Rider {
	return Rider{ID: id}
}
