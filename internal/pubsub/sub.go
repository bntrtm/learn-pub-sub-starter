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

type Acktype int

const (
	Ack Acktype = iota
	NackDiscard
	NackRequeue
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) Acktype,
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
			acktype := handler(data)
			switch acktype {
			case Ack:
				err = msg.Ack(false)
			case NackRequeue:
				err = msg.Nack(false, true)
			case NackDiscard:
				err = msg.Nack(false, false)
			default:
				log.Println("invalid acktype")
				continue
			}
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
