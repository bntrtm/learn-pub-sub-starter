package main

import (
	"fmt"

	"github.com/bntrtm/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bntrtm/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		prompt := "> "
		if ps.IsPaused {
			prompt = "||"
		}
		defer fmt.Print(prompt)
		gs.HandlePause(ps)
	}
}
