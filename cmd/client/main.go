package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bntrtm/learn-pub-sub-starter/internal/gamelogic"
	ps "github.com/bntrtm/learn-pub-sub-starter/internal/pubsub"
	"github.com/bntrtm/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	ConnString := "amqp://guest:guest@localhost:5672/"
	cxn, err := amqp.Dial(ConnString)
	if err != nil {
		log.Fatalf("could not open connection: %s", err)
	}
	defer cxn.Close()
	fmt.Println("Server connection  successful!")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not log in: %s", err)
	}

	_, _, err = ps.DeclareAndBind(
		cxn,
		"peril_direct",
		ps.BuildQueueName(routing.PauseKey, username),
		routing.PauseKey,
		ps.Transient)
	if err != nil {
		log.Fatalf("queue error: %s", err)
	}

	// client shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nPeril client shutting down...")
}
