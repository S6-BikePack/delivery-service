package main

import (
	"context"
	"delivery-service/config"
	"delivery-service/internal/core/services"
	"delivery-service/internal/handlers"
	"delivery-service/internal/repositories"
	"delivery-service/pkg/logging"
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

	logger, err := logging.NewSimpleLogger(cfg)

	if err != nil {
		panic(err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Database
	//--------------------------------------------------------------------------------------

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		logger.Warning(context.Background(), "Failed to connect to database: %s", err.Error())
	}

	if cfg.Database.Debug {
		db.Debug()
	}

	deliveryRepository, err := repositories.NewDeliveryRepository(db)

	if err != nil {
		logger.Warning(context.Background(), "Failed to create delivery repository: %s", err.Error())
	}

	riderRepository, err := repositories.NewRiderRepository(db)

	if err != nil {
		logger.Warning(context.Background(), "Failed to create rider repository: %s", err.Error())
	}

	serviceAreaRepository, err := repositories.NewServiceAreaRepository(db)

	if err != nil {
		logger.Warning(context.Background(), "Failed to create service area repository: %s", err.Error())
	}

	//--------------------------------------------------------------------------------------
	// Setup RabbitMQ
	//--------------------------------------------------------------------------------------

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		logger.Warning(context.Background(), "Failed to create RabbitMQ server: %s", err.Error())
	}

	rmqPublisher := services.NewRabbitMQPublisher(rmqServer, cfg)

	//--------------------------------------------------------------------------------------
	// Setup Services
	//--------------------------------------------------------------------------------------

	serviceAreaService := services.NewServiceAreaService(serviceAreaRepository)
	riderService := services.NewRiderService(riderRepository, rmqPublisher)
	routingService := services.NewMapboxService(cfg)
	deliveryService := services.NewDeliveryService(deliveryRepository, rmqPublisher, riderService, routingService)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, deliveryService, riderService, serviceAreaService, logger, cfg)

	//--------------------------------------------------------------------------------------
	// Setup HTTP server
	//--------------------------------------------------------------------------------------

	router := gin.New()

	deliveryHandler := handlers.NewRest(deliveryService, router, logger)
	deliveryHandler.SetupEndpoints()
	deliveryHandler.SetupSwagger()

	go rmqSubscriber.Listen()
	logger.Fatal(context.Background(), router.Run(cfg.Server.Port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	fmt.Println(returnValue)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
