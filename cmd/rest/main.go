package main

import (
	"delivery-service/config"
	"delivery-service/internal/core/services"
	"delivery-service/internal/core/services/logging_service"
	"delivery-service/internal/core/services/rabbitmq_service"
	"delivery-service/internal/handlers"
	"delivery-service/internal/repositories"
	"delivery-service/pkg/rabbitmq"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"

	"github.com/gin-gonic/gin"
)

const defaultConfig = "./config/local.config"

func main() {
	cfgPath := GetEnvOrDefault("config", defaultConfig)
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(err)
	}

	logger := logging_service.NewZerologLogger(cfg.Server.Service)

	//--------------------------------------------------------------------------------------
	// Setup Database
	//--------------------------------------------------------------------------------------

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		logger.Fatal(err)
	}

	if cfg.Database.Debug {
		db.Debug()
	}

	deliveryRepository, err := repositories.NewDeliveryRepository(db)

	if err != nil {
		logger.Fatal(err)
	}

	riderRepostory, err := repositories.NewRiderRepository(db)

	if err != nil {
		logger.Fatal(err)
	}

	serviceAreaRepository, err := repositories.NewServiceAreaRepository(db)

	if err != nil {
		logger.Fatal(err)
	}

	//--------------------------------------------------------------------------------------
	// Setup RabbitMQ
	//--------------------------------------------------------------------------------------

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		logger.Fatal(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer, cfg)

	//--------------------------------------------------------------------------------------
	// Setup Services
	//--------------------------------------------------------------------------------------

	serviceAreaService := services.NewServiceAreaService(serviceAreaRepository)
	riderService := services.NewRiderService(riderRepostory, rmqPublisher)
	deliveryService := services.NewDeliveryService(deliveryRepository, rmqPublisher, riderService)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, deliveryService, riderService, serviceAreaService, logger, cfg)

	//--------------------------------------------------------------------------------------
	// Setup HTTP server
	//--------------------------------------------------------------------------------------

	router := gin.New()

	deliveryHandler := handlers.NewRest(deliveryService, router, logger)
	deliveryHandler.SetupEndpoints()
	deliveryHandler.SetupSwagger()

	go rmqSubscriber.Listen()
	logger.Fatal(router.Run(cfg.Server.Port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	fmt.Println(returnValue)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
