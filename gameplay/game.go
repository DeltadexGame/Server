package gameplay

import (
	"deltadex/server/networking"
	"math/rand"
	"time"

	"github.com/Strum355/log"
	"github.com/google/uuid"
)

// GameStatus keeps track of what status the game is currently in
var GameStatus int = 0

// GameTurn is what turn the game is currently on
var GameTurn int = 0

// PlayerTurn decides who is currently taking their turn
var PlayerTurn int

// NextTurn declares if when the next player ends their turn if the game moves onto the next turn
var NextTurn bool

// PlayerOne is the first player
var PlayerOne *Player

// PlayerTwo is the second player
var PlayerTwo *Player

// GameID is used to identify each game
var GameID uuid.UUID

// Initialise initialises all of the variables in the game
func Initialise() {
	GameID = uuid.New()

	PlayerOne.Hand = []Card{}
	PlayerTwo.Hand = []Card{}
	PlayerOne.Deck = []Card{}
	PlayerTwo.Deck = []Card{}
	PlayerOne.Monsters = make([]Monster, 3)
	PlayerTwo.Monsters = make([]Monster, 3)

	PlayerOne.Energy = 10
	PlayerTwo.Energy = 10
	PlayerOne.Health = 50
	PlayerTwo.Health = 50

	GameTurn = 0
	GameStatus = 0
	PlayerTurn = 0
}

// Reset resets the game
func Reset() {
	Initialise()
}

// Start starts the game
func Start() {
	// Reset the game before starting
	Initialise()
	log.WithFields(log.Fields{
		"game_id": GameID,
	}).Info("Starting game, both players connected")

	// Set the game to the first turn and set status to in-progress
	GameTurn = 1
	GameStatus = 1

	// Decide the starting player
	rand.Seed(time.Now().UnixNano())
	PlayerTurn = rand.Intn(2) + 1
	log.WithFields(log.Fields{
		"starting_player": PlayerTurn,
		"player_one":      PlayerOne.Username,
		"player_two":      PlayerTwo.Username,
	}).Info("Starting player decided")

	// Send packets to the players to inform them who starts
	packetContent := map[string]interface{}{
		"starting_player": PlayerTurn,
		"energy":          PlayerOne.Energy,
		"health":          PlayerOne.MaxHealth,
		"player_one":      PlayerOne.Username,
		"player_two":      PlayerTwo.Username,
	}
	PlayerOne.SendPacket(networking.Packet{PacketID: networking.GameInitiationInformation, Content: packetContent})
	PlayerTwo.SendPacket(networking.Packet{PacketID: networking.GameInitiationInformation, Content: packetContent})

	// Send players packets with their starting hands
	card := Card{ID: 0, Name: "Zombie", Type: 0, Attack: 1, Health: 12, EnergyCost: 1, Ability: Ability{AbilityID: 0, Name: "Zombie", Description: "Monster resurrected at half health upon death", Targeted: false}}
	hand := []Card{card, card, card, card}

	packetContent = map[string]interface{}{
		"hand": hand,
	}
	PlayerOne.Hand = hand
	PlayerTwo.Hand = hand
	PlayerOne.SendPacket(networking.Packet{PacketID: networking.StartingHand, Content: packetContent})
	PlayerTwo.SendPacket(networking.Packet{PacketID: networking.StartingHand, Content: packetContent})
}

// EndTurn is run when each player ends their turn
func EndTurn(playerID int) {
	var player *Player
	if playerID == 1 {
		player = PlayerOne
	} else {
		player = PlayerTwo
	}
	for index := 0; index < len(player.Monsters); index++ {
		if player.Monsters[index] == (Monster{}) {
			continue
		}
		monster := player.Monsters[index]
		if player.OtherPlayer().Monsters[index] == (Monster{}) {
			player.OtherPlayer().Health -= monster.Attack
			// Check if they've died
			// Send packet
			continue
		}

		player.OtherPlayer().Monsters[index].Damage(monster.Attack)
		// Send packet
	}
}
