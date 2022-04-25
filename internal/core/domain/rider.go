package domain

type Rider struct {
	ID string `json:"id"`
}

func NewRider(id string) Rider {
	return Rider{ID: id}
}
