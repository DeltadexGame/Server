package gameplay

import (
	"deltadex/server/networking"
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

	packetContent := map[string]interface{}{
		"starting_player": PlayerTurn,
		"energy":          10,
	}

	PlayerOne.SendPacket(networking.Packet{PacketID: 2, Content: packetContent})
	PlayerTwo.SendPacket(networking.Packet{PacketID: 2, Content: packetContent})

	card := Card{ID: 0, Name: "Zombie", Type: 0, Attack: 1, Health: 12, EnergyCost: 1, Ability: Ability{AbilityID: 0, Name: "Zombie", Description: "Monster resurrected at half health upon death", Targeted: false}}
	hand := []Card{card, card, card, card}

	packetContent = map[string]interface{}{
		"hand": hand,
	}
	PlayerOne.SendPacket(networking.Packet{PacketID: 3, Content: packetContent})
	PlayerTwo.SendPacket(networking.Packet{PacketID: 3, Content: packetContent})

	for {
		time.Sleep(5 * time.Second)

		packetContent = map[string]interface{}{
			"hand": []Card{card},
		}
		PlayerOne.SendPacket(networking.Packet{PacketID: 3, Content: packetContent})
		PlayerTwo.SendPacket(networking.Packet{PacketID: 3, Content: packetContent})
	}
}
