package services

import (
	"delivery-service/config"
	"delivery-service/internal/core/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

type MapboxService struct {
	config *config.Config
}

func NewMapboxService(config *config.Config) *MapboxService {
	return &MapboxService{config: config}
}

func (s *MapboxService) GetRoute(pickup, destination domain.Location) (domain.Line, error) {
	url := fmt.Sprintf(
		"https://api.mapbox.com/directions/v5/mapbox/cycling/%f,%f;%f,%f?alternatives=false&continue_straight=true&geometries=geojson&overview=full&steps=false&access_token=%s",
		pickup.Longitude,
		pickup.Latitude,
		destination.Longitude,
		destination.Latitude,
		s.config.MapBox.AccessToken)

	resp, err := http.Get(url)

	if err != nil {
		return domain.Line{}, err
	}

	var body struct {
		Routes []struct {
			Geometry domain.Line
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return domain.Line{}, err
	}

	return body.Routes[0].Geometry, nil
}
