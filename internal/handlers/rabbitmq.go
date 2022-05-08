package handlers

import (
	"delivery-service/config"
	"delivery-service/internal/core/domain"
	"delivery-service/internal/core/interfaces"
	"delivery-service/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
)

type rabbitmqHandler struct {
	rabbitmq           *rabbitmq.RabbitMQ
	deliveryService    interfaces.DeliveryService
	riderService       interfaces.RiderService
	serviceAreaService interfaces.ServiceAreaService
	handlers           map[string]func(topic string, body []byte, handler *rabbitmqHandler) error
	logger             interfaces.LoggingService
	config             *config.Config
}

func NewRabbitMQ(rabbitmq *rabbitmq.RabbitMQ, deliveryService interfaces.DeliveryService, riderService interfaces.RiderService, serviceAreaService interfaces.ServiceAreaService, logger interfaces.LoggingService, config *config.Config) *rabbitmqHandler {
	return &rabbitmqHandler{
		rabbitmq:           rabbitmq,
		deliveryService:    deliveryService,
		riderService:       riderService,
		serviceAreaService: serviceAreaService,
		handlers: map[string]func(topic string, body []byte, handler *rabbitmqHandler) error{
			"rider.create": RiderCreate,
			"rider.update": RiderUpdate,
			"rider." + config.ServiceArea.Identifier + ".update.location": RiderUpdateLocation,
			"customer.create":         CustomerCreateOrUpdate,
			"customer.update.details": CustomerCreateOrUpdate,
			"parcel." + config.ServiceArea.Identifier + ".create": ParcelCreateOrUpdate,
			"service_area.create": ServiceAreaCreateOrUpdate,
		},
		logger: logger,
		config: config,
	}
}

func ServiceAreaCreateOrUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var serviceArea domain.ServiceArea
	if err := json.Unmarshal(body, &serviceArea); err != nil {
		return err
	}

	if err := handler.serviceAreaService.SaveOrUpdateServiceArea(serviceArea); err != nil {
		return err
	}

	return nil
}

func RiderCreate(topic string, body []byte, handler *rabbitmqHandler) error {
	var riderObject = struct {
		UserID string
		User   struct {
			Name string
		}
		Status      int
		ServiceArea domain.ServiceArea
		Capacity    domain.Dimensions
		Location    domain.Location
	}{}

	if err := json.Unmarshal(body, &riderObject); err != nil {
		return err
	}

	fmt.Println(riderObject)

	if _, err := handler.riderService.Create(riderObject.UserID, riderObject.User.Name, riderObject.ServiceArea.ID); err != nil {
		return err
	}

	return nil
}

func RiderUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var riderObject = struct {
		UserID string
		User   struct {
			Name string
		}
		Status      int
		ServiceArea domain.ServiceArea
		Capacity    domain.Dimensions
	}{}

	if err := json.Unmarshal(body, &riderObject); err != nil {
		return err
	}

	rider := domain.Rider{
		ID:            riderObject.UserID,
		Name:          riderObject.User.Name,
		ServiceAreaID: riderObject.ServiceArea.ID,
		IsActive:      riderObject.ServiceArea.ID == handler.config.ServiceArea.Id && riderObject.Status == 2,
	}

	if _, err := handler.riderService.Update(rider); err != nil {
		return err
	}

	return nil
}

func RiderUpdateLocation(topic string, body []byte, handler *rabbitmqHandler) error {
	message := struct {
		Id       string
		Location domain.Location
	}{}

	if err := json.Unmarshal(body, &message); err != nil {
		return err
	}

	if message.Id == "" {
		return nil
	}

	if err := handler.riderService.UpdateLocation(message.Id, message.Location); err != nil {
		return err
	}

	return nil
}

func CustomerCreateOrUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var customer domain.Customer

	if err := json.Unmarshal(body, &customer); err != nil {
		return err
	}

	if err := handler.deliveryService.SaveOrUpdateCustomer(customer); err != nil {
		return err
	}

	return nil
}

func ParcelCreateOrUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var parcel domain.Parcel

	if err := json.Unmarshal(body, &parcel); err != nil {
		return err
	}

	if err := handler.deliveryService.SaveOrUpdateParcel(parcel); err != nil {
		return err
	}

	return nil
}

func (handler *rabbitmqHandler) Listen() {

	q, err := handler.rabbitmq.Channel.QueueDeclare(
		handler.config.Server.Service+"-"+handler.config.ServiceArea.Identifier,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		handler.logger.Error(err)
	}

	for _, s := range maps.Keys(handler.handlers) {
		err = handler.rabbitmq.Channel.QueueBind(
			q.Name,
			s,
			handler.config.RabbitMQ.Exchange,
			false,
			nil)
		if err != nil {
			handler.logger.Error(err)
		}
	}

	msgs, err := handler.rabbitmq.Channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		handler.logger.Error(err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			fun, exist := handler.handlers[msg.RoutingKey]

			if exist {
				if err = fun(msg.RoutingKey, msg.Body, handler); err == nil {
					err = msg.Ack(false)
					if err != nil {
						handler.logger.Error(err)
					}

					continue
				}

				err = msg.Nack(false, true)
				if err != nil {
					handler.logger.Error(err)
				}

				continue
			}

			handler.logger.Warnf("No handler exists for %d", msg.RoutingKey)
			err = msg.Nack(false, true)
			if err != nil {
				handler.logger.Error(err)
			}
		}
	}()

	<-forever
}

type MessageHandler struct {
	topic   string
	handler func(topic string, body []byte, handler *rabbitmqHandler) error
}
