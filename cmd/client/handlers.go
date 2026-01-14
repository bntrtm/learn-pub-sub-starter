package main

import (
	"fmt"

	"github.com/bntrtm/learn-pub-sub-starter/internal/gamelogic"
	ps "github.com/bntrtm/learn-pub-sub-starter/internal/pubsub"
	"github.com/bntrtm/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) ps.Acktype {
	return func(playState routing.PlayingState) ps.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(playState)
		return ps.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.ArmyMove) ps.Acktype {
	return func(move gamelogic.ArmyMove) ps.Acktype {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		switch outcome {
		case gamelogic.MoveOutcomeMakeWar:
			err := ps.SendWarMessage(publishCh, move.Player.Username, gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			})
			if err != nil {
				return ps.NackRequeue
			}
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

func handlerRecognizeWar(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.RecognitionOfWar) ps.Acktype {
	return func(warRec gamelogic.RecognitionOfWar) ps.Acktype {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(warRec)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return ps.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return ps.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			err := ps.SendGameLog(publishCh,
				warRec.Attacker.Username,
				fmt.Sprintf("%s won a war against %s", winner, loser))
			if err != nil {
				return ps.NackRequeue
			}
			return ps.Ack
		case gamelogic.WarOutcomeYouWon:
			err := ps.SendGameLog(publishCh,
				warRec.Attacker.Username,
				fmt.Sprintf("%s won a war against %s", winner, loser))
			if err != nil {
				return ps.NackRequeue
			}
			return ps.Ack
		case gamelogic.WarOutcomeDraw:
			err := ps.SendGameLog(publishCh,
				warRec.Attacker.Username,
				fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser))
			if err != nil {
				return ps.NackRequeue
			}
			return ps.Ack
		default:
			fmt.Println("error: unrecognized war outcome")
			return ps.NackDiscard
		}
	}
}
