package main

import (
	"delivery-service/internal/core/services"
	"delivery-service/internal/core/services/logging_service"
	"delivery-service/internal/core/services/rabbitmq_service"
	"delivery-service/internal/handlers"
	"delivery-service/internal/repositories"
	"delivery-service/pkg/rabbitmq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"

	"github.com/gin-gonic/gin"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgres://user:password@localhost:5432/delivery"

func main() {
	logger := logging_service.NewZerologLogger("delivery-service")

	dbConn := GetEnvOrDefault("DATABASE", defaultDbConn)

	db, err := gorm.Open(postgres.Open(dbConn))
	db.Debug()

	if err != nil {
		logger.Fatal(err)
	}

	deliveryRepository, err := repositories.NewDeliveryRepository(db)

	if err != nil {
		logger.Fatal(err)
	}

	riderRepostory, err := repositories.NewRiderRepository(db)

	if err != nil {
		logger.Fatal(err)
	}

	rmqConn := GetEnvOrDefault("RABBITMQ", defaultRmqConn)

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		logger.Fatal(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer)

	riderService := services.NewRiderService(riderRepostory, rmqPublisher)
	deliveryService := services.NewDeliveryService(deliveryRepository, rmqPublisher, riderService)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, deliveryService, riderService, logger)

	router := gin.New()

	deliveryHandler := handlers.NewRest(deliveryService, router, logger)
	deliveryHandler.SetupEndpoints()
	deliveryHandler.SetupSwagger()

	port := GetEnvOrDefault("PORT", defaultPort)

	go rmqSubscriber.Listen("deliveryQueue")
	logger.Fatal(router.Run(port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
