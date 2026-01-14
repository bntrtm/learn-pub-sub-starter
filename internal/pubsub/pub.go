// Package pubsub offers helpful functions for defining queues
// for sending messaages via publishers, and consuming them
// via subscribers.
package pubsub

import (
	"context"
	"encoding/json"

	"github.com/bntrtm/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bntrtm/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishJSON publishes a JSON message to an AMQP exchange.
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

func SendPauseMessage(channel *amqp.Channel, isPaused bool) error {
	return PublishJSON(
		channel,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: isPaused},
	)
}

func SendMoveMessage(channel *amqp.Channel, username string, move gamelogic.ArmyMove) error {
	return PublishJSON(
		channel,
		routing.ExchangePerilTopic,
		RPattern(routing.ArmyMovesPrefix, "username"),
		move,
	)
}

func SendWarMessage(channel *amqp.Channel, username string, warRec gamelogic.RecognitionOfWar) error {
	return PublishJSON(
		channel,
		routing.ExchangePerilTopic,
		RPattern(routing.WarRecognitionsPrefix, username),
		warRec,
	)
}
