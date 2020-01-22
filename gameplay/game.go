package gameplay

import (
	"math/rand"
	"time"

	"github.com/Strum355/log"
)

// GameStatus keeps track of what status the game is currently in
var GameStatus int = 0

// GameTurn is what turn the game is currently on
var GameTurn int = 0

// PlayerTurn decides who is currently taking their turn
var PlayerTurn int

// PlayerOne is the first player
var PlayerOne *Player

// PlayerTwo is the second player
var PlayerTwo *Player

// Start starts the game
func Start() {
	log.Info("Starting game, both players connected")
	GameTurn = 1
	rand.Seed(time.Now().UnixNano())
	PlayerTurn = rand.Intn(2) + 1
	log.WithFields(log.Fields{
		"starting_player": PlayerTurn,
	}).Info("Starting player decided")
}
