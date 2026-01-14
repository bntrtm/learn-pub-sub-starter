package main

import (
	"fmt"
	"log"

	"github.com/bntrtm/learn-pub-sub-starter/internal/gamelogic"
	ps "github.com/bntrtm/learn-pub-sub-starter/internal/pubsub"
	"github.com/bntrtm/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()
	ConnString := "amqp://guest:guest@localhost:5672/"
	cxn, err := amqp.Dial(ConnString)
	if err != nil {
		log.Fatalf("could not open connection to RabbitMQ: %v", err)
	}
	defer cxn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	pub := ps.Publisher{}
	pauseChannel, err := cxn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	_, _, err = ps.DeclareAndBind(
		cxn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		ps.BuildQueueString(routing.GameLogSlug, "*"),
		ps.Durable)
	if err != nil {
		log.Fatalf("queue error: %v", err)
	}

GameLoop:
	for {
		words := gamelogic.GetInput("")
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Sending pause message...")
			err := pub.SendPauseMessage(pauseChannel, true)
			if err != nil {
				log.Println("Message send error: ", err)
			}
		case "resume":
			fmt.Println("Sending resume message...")
			err := pub.SendPauseMessage(pauseChannel, false)
			if err != nil {
				log.Println("Message send error: ", err)
			}
		case "quit":
			log.Println("Exiting...")
			break GameLoop
		default:
			fmt.Printf("Command '%s' not recognized\n", words[0])
		}
	}
}
