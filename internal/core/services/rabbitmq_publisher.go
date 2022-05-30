package services

import (
	"delivery-service/config"
	"delivery-service/internal/core/domain"
	"delivery-service/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmqPublisher struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Config     *config.Config
}

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ, cfg *config.Config) *rabbitmqPublisher {
	return &rabbitmqPublisher{Connection: rabbitmq.Connection, Channel: rabbitmq.Channel, Config: cfg}
}

func (rmq *rabbitmqPublisher) CreateDelivery(delivery domain.Delivery) error {
	err := rmq.publishJson("create", delivery)

	return err
}

func (rmq *rabbitmqPublisher) UpdateDelivery(delivery domain.Delivery) error {
	err := rmq.publishJson("update", delivery)

	return err
}

func (rmq *rabbitmqPublisher) StartDelivery(delivery domain.Delivery) error {
	body := struct{ ID string }{
		ID: delivery.ID,
	}

	err := rmq.publishJson("startDelivery", body)

	return err
}

func (rmq *rabbitmqPublisher) CompleteDelivery(delivery domain.Delivery) error {
	body := struct{ ID string }{
		ID: delivery.ID,
	}

	err := rmq.publishJson("completeDelivery", body)

	return err
}

func (rmq *rabbitmqPublisher) publishJson(topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = rmq.Channel.Publish(
		rmq.Config.RabbitMQ.Exchange,
		fmt.Sprintf("delivery.%s.%s", rmq.Config.ServiceArea.Identifier, topic),
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
