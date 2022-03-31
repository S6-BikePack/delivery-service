package main

import (
	"delivery-service/internal/core/services/delivery_service"
	"delivery-service/internal/core/services/rabbitmq_service"
	"delivery-service/internal/handlers"
	"delivery-service/internal/repositories/delivery_repository"
	"delivery-service/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

const defaultPort = ":1235"
const defaultRmqConn = "amqp://user:password@localhost:5672/"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := gin.Default()

	rmqConn := os.Getenv("RABBITMQ")
	if rmqConn == "" {
		rmqConn = defaultRmqConn
	}

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer)

	deliveryRepository, err := delivery_repository.NewCockroachDB("postgresql://root@localhost:26257/deliveries?sslmode=disable")

	if err != nil {
		panic(err)
	}

	deliveryService := delivery_service.New(deliveryRepository, rmqPublisher)

	handlers.NewGraphQL(router, deliveryService)
	rmqHandler := handlers.NewRabbitMQ(rmqServer, deliveryService)

	go rmqHandler.Listen("rider.#")
	log.Fatal(router.Run(port))
}
