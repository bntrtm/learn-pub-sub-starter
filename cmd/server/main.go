package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bntrtm/learn-pub-sub-starter/internal/pubsub"
	"github.com/bntrtm/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	ConnString := "amqp://guest:guest@localhost:5672/"
	cxn, err := amqp.Dial(ConnString)
	if err != nil {
		log.Fatalf("could not open connection: %s", err.Error())
	}
	defer cxn.Close()
	fmt.Println("Server start successful!")

	pauseChannel, err := cxn.Channel()
	if err != nil {
		panic(err)
	}
	err = pubsub.PublishJSON(pauseChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		panic(err)
	}

	// server shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nPeril server shutting down...")
}
