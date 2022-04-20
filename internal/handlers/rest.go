package handlers

import (
	"delivery-service/internal/core/domain"
	"delivery-service/internal/core/ports"
	"delivery-service/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"strconv"
	"strings"
	"time"

	"delivery-service/docs"
	_ "delivery-service/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type restHandler struct {
	deliveryService ports.DeliveryService
	router          *gin.Engine
	logger          ports.LoggingService
}

func NewRest(deliveryService ports.DeliveryService, router *gin.Engine, logger ports.LoggingService) *restHandler {
	return &restHandler{
		deliveryService: deliveryService,
		router:          router,
		logger:          logger,
	}
}

func (handler *restHandler) SetupEndpoints() {
	api := handler.router.Group("/api")
	api.GET("/deliveries", handler.GetAll)
	api.GET("/deliveries/:id", handler.Get)
	api.GET("/deliveries/radius/:latlon", handler.GetByDistance)
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
	delivery, err := handler.deliveryService.Get(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, dto.BuildResponseCreateDelivery(delivery))
}

// GetByDistance godoc
// @Summary  get delivery by distance
// @Schemes
// @Param        latlon     path  string           true  "Latitude,Longitude"
// @Param        radius    query     int  false  "radius of search in meters (default = 1000)"
// @Description  gets a delivery from the system based on the distance to the given point
// @Produce      json
// @Success      200  {object}  domain.Delivery
// @Router       /api/deliveries/radius/{latlon} [get]
func (handler *restHandler) GetByDistance(c *gin.Context) {
	latlon := strings.Split(c.Param("latlon"), ",")

	lat, err := strconv.ParseFloat(latlon[0], 64)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	lon, err := strconv.ParseFloat(latlon[1], 64)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	radiusQuery := c.Query("radius")
	var radius int64

	if radiusQuery != "" {
		radius, err = strconv.ParseInt(radiusQuery, 10, 32)
	} else {
		radius = 1000
	}

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	location := domain.Location{
		Latitude:  lat,
		Longitude: lon,
	}

	deliveries := handler.deliveryService.GetByDistance(location, int(radius))

	c.JSON(200, deliveries)
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
		handler.logger.Error(err)
		c.AbortWithStatus(500)
	}

	pickupTime := time.Unix(body.PickupTime, 0)

	delivery, err := handler.deliveryService.Create(body.ParcelId, body.OwnerId, body.PickupPoint, body.DeliveryPoint, pickupTime)

	if err != nil {
		handler.logger.Error(err)
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, delivery)
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
		handler.logger.Error(err)
		c.AbortWithStatus(500)
	}

	delivery, err := handler.deliveryService.AssignRider(c.Param("id"), body.RiderId)

	if err != nil {
		handler.logger.Error(err)
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
	delivery, err := handler.deliveryService.StartDelivery(c.Param("id"))

	if err != nil {
		handler.logger.Error(err)
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
	delivery, err := handler.deliveryService.CompleteDelivery(c.Param("id"))

	if err != nil {
		handler.logger.Error(err)
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, delivery)
}
