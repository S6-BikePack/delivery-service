package handlers

import (
	"delivery-service/internal/core/ports"
	"delivery-service/pkg/dto"
	"github.com/google/uuid"
	time2 "time"
)
import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	deliveryService ports.DeliveryService
}

func NewRest(deliveryService ports.DeliveryService) *HTTPHandler {
	return &HTTPHandler{
		deliveryService: deliveryService,
	}
}

// GetAll godoc
// @Summary  get all deliveries
// @Schemes
// @Description  gets all deliveries in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  []domain.Delivery
// @Router       /api/deliveries [get]
func (hdl *HTTPHandler) GetAll(c *gin.Context) {
	deliveries, err := hdl.deliveryService.GetAll()

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, deliveries)
}

// Get godoc
// @Summary  get rider
// @Schemes
// @Param        id     path  string           true  "Rider id"
// @Description  gets a rider from the system by its ID
// @Produce      json
// @Success      200  {object}  domain.Delivery
// @Router       /api/deliveries/{id} [get]
func (hdl *HTTPHandler) Get(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	delivery, err := hdl.deliveryService.Get(uid)

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, delivery)
}

// Create godoc
// @Summary  create rider
// @Schemes
// @Description  creates a new rider
// @Accept       json
// @Param        rider  body  dto.BodyCreateDelivery  true  "Add rider"
// @Produce      json
// @Success      200  {object}  dto.ResponseCreateDelivery
// @Router       /api/deliveries [post]
func (hdl *HTTPHandler) Create(c *gin.Context) {
	body := dto.BodyCreateDelivery{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	parcelId, err := uuid.Parse(body.ParcelId)

	if err != nil {
		c.AbortWithStatus(409)
	}

	time := time2.Unix(body.PickupTime, 0)

	delivery, err := hdl.deliveryService.Create(parcelId, body.PickupPoint, body.DeliveryPoint, time)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseCreateDelivery(delivery))
}

//func (hdl *HTTPHandler) UpdateRider(c *gin.Context) {
//	body := BodyUpdate{}
//	err := c.BindJSON(&body)
//
//	if err != nil {
//		c.AbortWithStatus(500)
//	}
//
//	uid, err := uuid.Parse(c.Param("id"))
//
//	if err != nil {
//		c.AbortWithStatus(400)
//		return
//	}
//
//	delivery, err := hdl.deliveryService.AssignRider(uid, body.Name, body.Status)
//
//	if err != nil {
//		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
//		return
//	}
//
//	c.JSON(200, BuildResponseUpdate(delivery))
//}
