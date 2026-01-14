package main

import (
	"fmt"

	"github.com/bntrtm/learn-pub-sub-starter/internal/gamelogic"
	ps "github.com/bntrtm/learn-pub-sub-starter/internal/pubsub"
	"github.com/bntrtm/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) ps.Acktype {
	return func(playState routing.PlayingState) ps.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(playState)
		return ps.Ack
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) ps.Acktype {
	return func(am gamelogic.ArmyMove) ps.Acktype {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(am)
		switch outcome {
		case gamelogic.MoveOutcomeMakeWar:
			return ps.Ack
		case gamelogic.MoveOutComeSafe:
			return ps.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return ps.NackDiscard
		}
		fmt.Println("error: unrecognized move outcome")
		return ps.NackDiscard
	}
}
