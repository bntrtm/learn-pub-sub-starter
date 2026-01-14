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
	fmt.Println("Starting Peril client...")
	ConnString := "amqp://guest:guest@localhost:5672/"
	cxn, err := amqp.Dial(ConnString)
	if err != nil {
		log.Fatalf("could not open connection: %v", err)
	}
	defer cxn.Close()
	fmt.Println("Server connection  successful!")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not log in: %v", err)
	}

	_, _, err = ps.DeclareAndBind(
		cxn,
		routing.ExchangePerilDirect,
		ps.BuildQueueName(routing.PauseKey, username),
		routing.PauseKey,
		ps.Transient)
	if err != nil {
		log.Fatalf("queue error: %v", err)
	}

	gameState := gamelogic.NewGameState(username)

REPL:
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			err := gameState.CommandSpawn(words)
			if err != nil {
				log.Println(err)
			}
		case "move":
			move, err := gameState.CommandMove(words)
			if err != nil {
				log.Println(err)
			}
			log.Printf("moved unit %d to %s", move.Units[0].ID, move.ToLocation)
		case "status":
			gameState.CommandStatus()
		case "spam":
			fmt.Println("Spamming not implemented yet!")
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			gamelogic.PrintQuit()
			break REPL
		default:
			fmt.Printf("Command '%s' not recognized\n", words[0])
		}
	}
}
