// Package pubsub offers helpful functions for defining queues
// and sending JSON-based messages
package pubsub

import (
	"context"
	"encoding/json"

	"github.com/bntrtm/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishJSON publishes a JSON message to an AMQP exchange.z
func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/JSON",
		Body:        jsonBytes,
	})
	if err != nil {
		return err
	}
	return nil
}

type Publisher struct{}

func (p *Publisher) SendPauseMessage(channel *amqp.Channel, isPaused bool) error {
	return PublishJSON(
		channel,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: isPaused},
	)
}
