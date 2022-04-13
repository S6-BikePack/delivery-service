package main

import (
	"delivery-service/internal/core/services/delivery_service"
	"delivery-service/internal/core/services/rabbitmq_service"
	"delivery-service/internal/handlers"
	"delivery-service/internal/repositories/delivery_repository"
	"delivery-service/pkg/rabbitmq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgresql://user:password@localhost:5432/delivery"

func main() {
	dbConn := GetEnvOrDefault("DATABASE", defaultDbConn)

	db, err := gorm.Open(postgres.Open(dbConn))

	if err != nil {
		panic(err)
	}

	deliveryRepository, err := delivery_repository.NewCockroachDB(db)

	if err != nil {
		panic(err)
	}

	rmqConn := GetEnvOrDefault("RABBITMQ", defaultRmqConn)

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer)

	deliveryService := delivery_service.New(deliveryRepository, rmqPublisher)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, deliveryService)

	router := gin.New()

	deliveryHandler := handlers.NewRest(deliveryService, router)
	deliveryHandler.SetupEndpoints()
	deliveryHandler.SetupSwagger()

	port := GetEnvOrDefault("PORT", defaultPort)

	go rmqSubscriber.Listen("deliveryQueue")
	log.Fatal(router.Run(port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
