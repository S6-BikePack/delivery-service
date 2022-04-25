package handlers

import (
	"delivery-service/internal/core/domain"
	"delivery-service/internal/core/interfaces"
	"delivery-service/pkg/rabbitmq"
	"encoding/json"
	"golang.org/x/exp/maps"
)

type rabbitmqHandler struct {
	rabbitmq        *rabbitmq.RabbitMQ
	deliveryService interfaces.DeliveryService
	riderService    interfaces.RiderService
	handlers        map[string]func(topic string, body []byte, handler *rabbitmqHandler) error
	logger          interfaces.LoggingService
}

func NewRabbitMQ(rabbitmq *rabbitmq.RabbitMQ, deliveryService interfaces.DeliveryService, riderService interfaces.RiderService, logger interfaces.LoggingService) *rabbitmqHandler {
	return &rabbitmqHandler{
		rabbitmq:        rabbitmq,
		deliveryService: deliveryService,
		riderService:    riderService,
		handlers: map[string]func(topic string, body []byte, handler *rabbitmqHandler) error{
			"rider.create":            RiderCreate,
			"rider.update":            RiderUpdate,
			"rider.update.location":   RiderUpdateLocation,
			"customer.create":         CustomerCreateOrUpdate,
			"customer.update.details": CustomerCreateOrUpdate,
			"parcel.create":           ParcelCreateOrUpdate,
		},
		logger: logger,
	}
}

func RiderCreate(topic string, body []byte, handler *rabbitmqHandler) error {
	var rider domain.Rider

	if err := json.Unmarshal(body, &rider); err != nil {
		return err
	}

	if _, err := handler.riderService.Create(rider.ID, rider.Name, rider.ServiceArea); err != nil {
		return err
	}

	return nil
}

func RiderUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var rider domain.Rider

	if err := json.Unmarshal(body, &rider); err != nil {
		return err
	}

	if _, err := handler.riderService.Update(rider); err != nil {
		return err
	}

	return nil
}

func RiderUpdateLocation(topic string, body []byte, handler *rabbitmqHandler) error {
	message := struct {
		id       string
		location domain.Location
	}{}

	if err := json.Unmarshal(body, &message); err != nil {
		return err
	}

	if err := handler.riderService.UpdateLocation(message.id, message.location); err != nil {
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

func (handler *rabbitmqHandler) Listen(queue string) {

	q, err := handler.rabbitmq.Channel.QueueDeclare(
		queue,
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
			"topics",
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
					msg.Ack(false)

					continue
				}

				handler.logger.Error(err)
				msg.Nack(false, true)

				continue
			}

			handler.logger.Warnf("No handler exists for %d", msg.RoutingKey)
			msg.Nack(false, true)
		}
	}()

	<-forever
}

type MessageHandler struct {
	topic   string
	handler func(topic string, body []byte, handler *rabbitmqHandler) error
}
