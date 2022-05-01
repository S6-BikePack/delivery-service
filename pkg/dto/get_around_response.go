package dto

type GetAroundResponse struct {
	Deliveries DeliveryListResponse `json:"deliveries"`
	Radius     int                  `json:"radius"`
}

func CreateGetAroundResponse(deliveries DeliveryListResponse, radius int) GetAroundResponse {
	return GetAroundResponse{
		Deliveries: deliveries,
		Radius:     radius,
	}
}
