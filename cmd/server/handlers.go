package main

import (
	"fmt"

	"github.com/bntrtm/learn-pub-sub-starter/internal/gamelogic"
	ps "github.com/bntrtm/learn-pub-sub-starter/internal/pubsub"
	"github.com/bntrtm/learn-pub-sub-starter/internal/routing"
)

func handlerLogs() func(routing.GameLog) ps.Acktype {
	return func(gameLog routing.GameLog) ps.Acktype {
		defer fmt.Print("> ")
		err := gamelogic.WriteLog(gameLog)
		if err != nil {
			return ps.NackDiscard
		}
		return ps.Ack
	}
}
