package handlers

import (
	"delivery-service/internal/core/ports"
	"delivery-service/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"time"

	"delivery-service/docs"
	_ "delivery-service/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type restHandler struct {
	deliveryService ports.DeliveryService
	router          *gin.Engine
}

func NewRest(deliveryService ports.DeliveryService, router *gin.Engine) *restHandler {
	return &restHandler{
		deliveryService: deliveryService,
		router:          router,
	}
}

func (handler *restHandler) SetupEndpoints() {
	api := handler.router.Group("/api")
	api.GET("/deliveries", handler.GetAll)
	api.GET("/deliveries/:id", handler.Get)
	api.POST("/deliveries", handler.Create)
	api.POST("/deliveries/:id/rider", handler.AssignRider)
	api.GET("/deliveries/:id/start", handler.StartDelivery)
	api.GET("/deliveries/:id/complete", handler.CompleteDelivery)
}

func (handler *restHandler) SetupSwagger() {
	docs.SwaggerInfo.Title = "Delivery service API"
	docs.SwaggerInfo.Description = "The delivery service manages all deliveries for the BikePack system."

	handler.router.GET("/delivery-service/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// GetAll godoc
// @Summary  get all deliveries
// @Schemes
// @Description  gets all deliveries in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  []domain.Delivery
// @Router       /api/deliveries [get]
func (handler *restHandler) GetAll(c *gin.Context) {
	deliveries, err := handler.deliveryService.GetAll()

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, deliveries)
}

// Get godoc
// @Summary  get delivery
// @Schemes
// @Param        id     path  string           true  "Delivery id"
// @Description  gets a delivery from the system by its ID
// @Produce      json
// @Success      200  {object}  domain.Delivery
// @Router       /api/deliveries/{id} [get]
func (handler *restHandler) Get(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	delivery, err := handler.deliveryService.Get(uid)

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, delivery)
}

// Create godoc
// @Summary  create delivery
// @Schemes
// @Description  creates a new delivery
// @Accept       json
// @Param        rider  body  dto.BodyCreateDelivery  true  "Add delivery"
// @Produce      json
// @Success      200  {object}  dto.ResponseCreateDelivery
// @Router       /api/deliveries [post]
func (handler *restHandler) Create(c *gin.Context) {
	body := dto.BodyCreateDelivery{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	parcelId, err := uuid.Parse(body.ParcelId)

	if err != nil {
		c.AbortWithStatus(409)
	}

	pickupTime := time.Unix(body.PickupTime, 0)

	delivery, err := handler.deliveryService.Create(parcelId, body.PickupPoint, body.DeliveryPoint, pickupTime)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseCreateDelivery(delivery))
}

// AssignRider godoc
// @Summary  assign rider
// @Schemes
// @Description  assigns a rider to a delivery
// @Accept       json
// @Param        delivery  body  dto.BodyAssignRider  true  "Assign rider"
// @Produce      json
// @Success      200  {object}  dto.ResponseAssignRider
// @Router       /api/deliveries/{id}/rider [post]
func (handler *restHandler) AssignRider(c *gin.Context) {
	body := dto.BodyAssignRider{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	riderId, err := uuid.Parse(body.RiderId)

	if err != nil {
		c.AbortWithStatus(409)
	}

	delivery, err := handler.deliveryService.AssignRider(uid, riderId)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseCreateDelivery(delivery))
}

// StartDelivery godoc
// @Summary  start delivery
// @Schemes
// @Description  starts a delivery
// @Produce      json
// @Success      200  {object}  domain.Delivery
// @Router       /api/deliveries/{id}/start [get]
func (handler *restHandler) StartDelivery(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	delivery, err := handler.deliveryService.StartDelivery(uid)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, delivery)
}

// CompleteDelivery godoc
// @Summary  complete delivery
// @Schemes
// @Description  completes a delivery
// @Produce      json
// @Success      200  {object}  domain.Delivery
// @Router       /api/deliveries/{id}/complete [get]
func (handler *restHandler) CompleteDelivery(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	delivery, err := handler.deliveryService.CompleteDelivery(uid)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, delivery)
}
