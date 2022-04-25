package rabbitmq_service

import (
	"delivery-service/internal/core/domain"
	"delivery-service/pkg/rabbitmq"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmqPublisher rabbitmq.RabbitMQ

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ) *rabbitmqPublisher {
	return &rabbitmqPublisher{Connection: rabbitmq.Connection, Channel: rabbitmq.Channel}
}

func (rmq *rabbitmqPublisher) CreateDelivery(delivery domain.Delivery) error {
	err := rmq.publishJson("delivery.create", delivery)

	return err
}

func (rmq *rabbitmqPublisher) UpdateDelivery(delivery domain.Delivery) error {
	err := rmq.publishJson("delivery.update", delivery)

	return err
}

func (rmq *rabbitmqPublisher) StartDelivery(delivery domain.Delivery) error {
	body := struct{ ID string }{
		ID: delivery.ID,
	}

	err := rmq.publishJson("delivery.startDelivery", body)

	return err
}

func (rmq *rabbitmqPublisher) CompleteDelivery(delivery domain.Delivery) error {
	body := struct{ ID string }{
		ID: delivery.ID,
	}

	err := rmq.publishJson("delivery.completeDelivery", body)

	return err
}

func (rmq *rabbitmqPublisher) publishJson(topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = rmq.Channel.Publish(
		"topics",
		topic,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         js,
		},
	)

	return err
}
