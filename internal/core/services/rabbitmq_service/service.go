package rabbitmq_service

import (
	"delivery-service/internal/core/domain"
	"delivery-service/pkg/rabbitmq"
	"encoding/json"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmqPublisher rabbitmq.RabbitMQ

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ) *rabbitmqPublisher {
	return &rabbitmqPublisher{Connection: rabbitmq.Connection, Channel: rabbitmq.Channel}
}

func (rmq *rabbitmqPublisher) CreateDelivery(delivery domain.Delivery) error {
	js, err := json.Marshal(delivery)

	if err != nil {
		return err
	}

	err = rmq.publishMessage("delivery.create", js)

	return err
}

func (rmq *rabbitmqPublisher) UpdateDelivery(delivery domain.Delivery) error {
	js, err := json.Marshal(delivery)

	if err != nil {
		return err
	}

	err = rmq.publishMessage("delivery.update", js)

	return err
}

func (rmq *rabbitmqPublisher) StartDelivery(delivery domain.Delivery) error {
	js, err := json.Marshal(struct{ ID uuid.UUID }{
		ID: delivery.ID,
	})

	if err != nil {
		return err
	}

	err = rmq.publishMessage("delivery.startDelivery", js)

	return err
}

func (rmq *rabbitmqPublisher) CompleteDelivery(delivery domain.Delivery) error {
	js, err := json.Marshal(struct{ ID uuid.UUID }{
		ID: delivery.ID,
	})

	if err != nil {
		return err
	}

	err = rmq.publishMessage("delivery.completeDelivery", js)

	return err
}

func (rmq *rabbitmqPublisher) publishMessage(key string, body []byte) error {
	err := rmq.Channel.Publish(
		"topics",
		key,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)

	return err
}

func (rmq *rabbitmqPublisher) GetRider(riderId uuid.UUID) (domain.Rider, error) {
	q, err := rmq.Channel.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)

	defer rmq.Channel.QueueDelete(q.Name, true, false, false)

	if err != nil {
		return domain.Rider{}, err
	}

	msgs, err := rmq.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return domain.Rider{}, err
	}

	corrId := uuid.New().String()

	err = rmq.Channel.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(riderId.String()),
		})

	if err != nil {
		return domain.Rider{}, err
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			var rider domain.Rider

			if err = json.Unmarshal(d.Body, &rider); err != nil {
				return domain.Rider{}, err
			}

			return rider, nil
		}
	}

	return domain.Rider{}, nil
}
