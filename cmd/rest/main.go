package main

import (
	"delivery-service/internal/core/services/delivery_service"
	"delivery-service/internal/core/services/rabbitmq_service"
	"delivery-service/internal/handlers"
	"delivery-service/internal/repositories/delivery_repository"
	"delivery-service/pkg/rabbitmq"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"delivery-service/docs"
	_ "delivery-service/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgresql://user:password@localhost:5432/delivery"

func main() {
	setupSwagger()

	dbConn := os.Getenv("DATABASE")
	if dbConn == "" {
		dbConn = defaultDbConn
	}

	deliveryRepository, err := delivery_repository.NewCockroachDB(dbConn)

	if err != nil {
		panic(err)
	}

	rmqConn := os.Getenv("RABBITMQ")
	if rmqConn == "" {
		rmqConn = defaultRmqConn
	}

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer)

	deliveryService := delivery_service.New(deliveryRepository, rmqPublisher)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, deliveryService)

	deliveryHandler := handlers.NewRest(deliveryService)

	router := gin.New()

	api := router.Group("/api")

	api.GET("/deliveries", deliveryHandler.GetAll)
	api.GET("/deliveries/:id", deliveryHandler.Get)
	api.POST("/deliveries", deliveryHandler.Create)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	go rmqSubscriber.Listen("rider.#")
	log.Fatal(router.Run(port))
}

func setupSwagger() {
	docs.SwaggerInfo.Title = "Delivery service API"
	docs.SwaggerInfo.Description = "The delivery service manages all deliveries for the BikePack system."
}
