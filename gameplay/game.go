package gameplay

import (
	"deltadex/gameplay/events"
	"deltadex/server/networking"
	"math/rand"
	"reflect"
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
	custom := make(map[string]map[string]reflect.Value)
	custom["deltadex/gameplay"] = make(map[string]reflect.Value)
	custom["deltadex/gameplay"]["Ability"] = reflect.ValueOf(Ability{})
	custom["deltadex/gameplay"]["Monster"] = reflect.ValueOf(Monster{})
	custom["deltadex/gameplay"]["Card"] = reflect.ValueOf(Card{})
	custom["deltadex/gameplay/events"] = make(map[string]reflect.Value)
	custom["deltadex/gameplay/events"]["Event"] = reflect.ValueOf(events.Event{})
	custom["deltadex/gameplay/events"]["EventID"] = reflect.ValueOf(events.MonsterAttackEvent)
	events.LoadScripts(custom)

	GameID = uuid.New()

	PlayerOne.Hand = []Card{}
	PlayerTwo.Hand = []Card{}
	PlayerOne.Deck = []Card{}
	PlayerTwo.Deck = []Card{}
	PlayerOne.Monsters = make([]Monster, 3)
	PlayerTwo.Monsters = make([]Monster, 3)

	PlayerOne.Energy = 10
	PlayerTwo.Energy = 10
	PlayerOne.MaxHealth = 50
	PlayerTwo.MaxHealth = 50
	PlayerOne.Health = PlayerOne.MaxHealth
	PlayerTwo.Health = PlayerTwo.MaxHealth

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
	selfContent := map[string]interface{}{
		"username": PlayerOne.Username,
		"id":       PlayerOne.ID,
		"energy":   PlayerOne.Energy,
		"health":   PlayerOne.Health,
		"starting": PlayerTurn == PlayerOne.ID,
	}
	PlayerOne.SendPacket(networking.Packet{PacketID: networking.SelfInitiationInformation, Content: selfContent})
	PlayerTwo.SendPacket(networking.Packet{PacketID: networking.OpponentInitiationInformation, Content: selfContent})

	selfContent = map[string]interface{}{
		"username": PlayerTwo.Username,
		"id":       PlayerTwo.ID,
		"energy":   PlayerTwo.Energy,
		"health":   PlayerTwo.Health,
		"starting": PlayerTurn == PlayerTwo.ID,
	}
	PlayerTwo.SendPacket(networking.Packet{PacketID: networking.SelfInitiationInformation, Content: selfContent})
	PlayerOne.SendPacket(networking.Packet{PacketID: networking.OpponentInitiationInformation, Content: selfContent})

	// Send players packets with their starting hands
	card := Card{ID: 0, Name: "Zombie", Type: 0, Attack: 2, Health: 4, EnergyCost: 2, Ability: Ability{AbilityID: 0, Name: "Zombie", Description: "Monster resurrected at half health upon death", Targeted: false}}
	hand := []Card{card, card, card, card}

	packetContent := map[string]interface{}{
		"hand": hand,
	}
	PlayerOne.Hand = hand
	PlayerTwo.Hand = hand
	PlayerOne.SendPacket(networking.Packet{PacketID: networking.StartingHand, Content: packetContent})
	PlayerTwo.SendPacket(networking.Packet{PacketID: networking.StartingHand, Content: packetContent})

	PlayerOne.SendPacket(networking.Packet{PacketID: networking.OpponentStartingHand, Content: map[string]interface{}{"hand": len(PlayerTwo.Hand)}})
	PlayerTwo.SendPacket(networking.Packet{PacketID: networking.OpponentStartingHand, Content: map[string]interface{}{"hand": len(PlayerOne.Hand)}})
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
			content := map[string]interface{}{
				"ownership": true,
				"position":  index,
				"died":      false,
				"health":    player.OtherPlayer().Health,
				"self":      false,
			}
			player.SendPacket(networking.Packet{PacketID: networking.EndTurnPlayerAttack, Content: content})
			content["self"] = true
			content["ownership"] = false
			player.OtherPlayer().SendPacket(networking.Packet{PacketID: networking.EndTurnPlayerAttack, Content: content})
			continue
		}

		player.OtherPlayer().Monsters[index].Damage(monster.Attack)
		died := false
		if player.OtherPlayer().Monsters[index].Health <= 0 {
			died = true
		}
		if died {
			if player.OtherPlayer().Monsters[index].Ability.AbilityID == 0 {
				player.OtherPlayer().Monsters[index].MaxHealth = player.OtherPlayer().Monsters[index].MaxHealth / 2
				player.OtherPlayer().Monsters[index].Health = player.OtherPlayer().Monsters[index].MaxHealth
				died = false
			}
		}
		content := map[string]interface{}{
			"ownership": false,
			"position":  index,
			"died":      died,
			"monster":   player.OtherPlayer().Monsters[index],
		}
		player.OtherPlayer().SendPacket(networking.Packet{PacketID: networking.EndTurnMonsterAttack, Content: content})
		content["ownership"] = true
		player.SendPacket(networking.Packet{PacketID: networking.EndTurnMonsterAttack, Content: content})

		if died {
			player.OtherPlayer().Monsters[index] = Monster{}
		}
	}
}
