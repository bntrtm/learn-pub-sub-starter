package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

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

// BuildQueueName concatenates two strings with an "."
// to produce an AMQP queue name.
func BuildQueueName(routingKey, name string) string {
	return routingKey + "." + name
}
