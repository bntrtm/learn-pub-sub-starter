package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	channel, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType)
	if err != nil {
		return err
	}
	deliveryChan, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range deliveryChan {
			var data T
			err := json.Unmarshal(msg.Body, &data)
			if err != nil {
				log.Println(err)
				continue
			}
			handler(data)
			err = msg.Ack(false)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	return nil
}

// DeclareAndBind creates and binds a queue to an AMQP exchange.
func DeclareAndBind(
	cxn *amqp.Connection,
	exchange,
	queueName,
	key string,

	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := cxn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	queue, err := channel.QueueDeclare(
		queueName,
		queueType == Durable,
		queueType == Transient,
		queueType == Transient,
		false,
		nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = channel.QueueBind(
		queue.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return channel, queue, nil
}

// BuildQueueString concatenates two strings with an "."
// to produce an AMQP queue string.
func BuildQueueString(routingKey, name string) string {
	return routingKey + "." + name
}
