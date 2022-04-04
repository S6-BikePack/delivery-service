package handlers

import (
	"delivery-service/internal/core/domain"
	"delivery-service/internal/core/ports"
	"delivery-service/pkg/rabbitmq"
	"encoding/json"
	"golang.org/x/exp/maps"
)

type rabbitmqHandler struct {
	rabbitmq *rabbitmq.RabbitMQ
	service  ports.DeliveryService
	handlers map[string]func(topic string, body []byte, handler *rabbitmqHandler) error
}

func NewRabbitMQ(rabbitmq *rabbitmq.RabbitMQ, service ports.DeliveryService) *rabbitmqHandler {
	return &rabbitmqHandler{
		rabbitmq: rabbitmq,
		service:  service,
		handlers: map[string]func(topic string, body []byte, handler *rabbitmqHandler) error{
			"rider.create":            RiderCreateOrUpdate,
			"rider.update":            RiderCreateOrUpdate,
			"customer.create":         CustomerCreateOrUpdate,
			"customer.update.details": CustomerCreateOrUpdate,
		},
	}
}

func RiderCreateOrUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var rider domain.Rider

	if err := json.Unmarshal(body, &rider); err != nil {
		return err
	}

	if err := handler.service.SaveOrUpdateRider(rider); err != nil {
		return err
	}

	return nil
}

func CustomerCreateOrUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var customer domain.Customer

	if err := json.Unmarshal(body, &customer); err != nil {
		return err
	}

	if err := handler.service.SaveOrUpdateCustomer(customer); err != nil {
		return err
	}

	return nil
}

func (handler *rabbitmqHandler) Listen() {

	q, err := handler.rabbitmq.Channel.QueueDeclare(
		"deliveryQueue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	for _, s := range maps.Keys(handler.handlers) {
		err = handler.rabbitmq.Channel.QueueBind(
			q.Name,
			s,
			"topics",
			false,
			nil)
		if err != nil {
			return
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
		panic(err)
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
			}

			msg.Nack(false, true)
		}
	}()

	<-forever
}

type MessageHandler struct {
	topic   string
	handler func(topic string, body []byte, handler *rabbitmqHandler) error
}
