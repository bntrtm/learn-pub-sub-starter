package main

import (
	"fmt"
	"log"
	"strconv"

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

	gameState := gamelogic.NewGameState(username)

	err = ps.SubscribeJSON(cxn,
		routing.ExchangePerilDirect,
		ps.RPattern(routing.PauseKey, username),
		routing.PauseKey,
		ps.Transient,
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatalf("queue error: %v", err)
	}

	pubChannel, err := cxn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	err = ps.SubscribeJSON(cxn,
		routing.ExchangePerilTopic,
		ps.RPattern(routing.ArmyMovesPrefix, username),
		ps.RPattern(routing.ArmyMovesPrefix, "*"),
		ps.Transient,
		handlerMove(gameState, pubChannel),
	)
	if err != nil {
		log.Fatalf("queue error: %v", err)
	}

	err = ps.SubscribeJSON(cxn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		ps.RPattern(routing.WarRecognitionsPrefix, "*"),
		ps.Durable,
		handlerRecognizeWar(gameState, pubChannel),
	)
	if err != nil {
		log.Fatalf("queue error: %v", err)
	}

REPL:
	for {
		words := gamelogic.GetInput("")
		if len(words) == 0 {
			continue
		}
		if gameState.Paused && words[0] != "quit" {
			words[0] = "GAME_PAUSED"
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
			err = ps.SendMoveMessage(
				pubChannel,
				username,
				move)
			if err != nil {
				log.Println(err)
			}
			log.Printf("moved %d units to %s", len(move.Units), move.ToLocation)
		case "status":
			gameState.CommandStatus()
		case "spam":
			if len(words) == 1 {
				fmt.Println("You can demonstrate pubsub backpressure by including a numeric value alongside the 'spam' command.")
				continue
			}
			n, err := strconv.ParseInt(words[1], 0, 64)
			if err != nil {
				log.Println("Usage: spam <number>")
			}
			for i := 0; i < int(n); i++ {
				malLog := gamelogic.GetMaliciousLog()
				err := ps.SendGameLog(pubChannel, username, malLog)
				if err != nil {
					log.Println(err)
				}
			}
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			gamelogic.PrintQuit()
			break REPL
		case "GAME_PAUSED":
			fmt.Println("|| The game is paused.")
			fmt.Println("|| Use 'quit' or wait for play to resume.")
			continue
		default:
			fmt.Printf("Command '%s' not recognized\n", words[0])
		}
	}
}
