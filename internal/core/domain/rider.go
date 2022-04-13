package domain

type Rider struct {
	ID   string
	Name string
}

func NewRider(id string, name string) Rider {
	return Rider{ID: id, Name: name}
}
