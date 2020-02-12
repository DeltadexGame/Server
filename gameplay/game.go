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

var (
	CurGame Game = Game{}
)

type Game struct {
	// GameStatus keeps track of what status the game is currently in
	GameStatus int

	// GameTurn is what turn the game is currently on
	GameTurn int

	// PlayerTurn decides who is currently taking their turn
	PlayerTurn int

	// NextTurn declares if when the next player ends their turn if the game moves onto the next turn
	NextTurn bool

	// PlayerOne is the first player
	PlayerOne *Player

	// PlayerTwo is the second player
	PlayerTwo *Player

	// GameID is used to identify each game
	GameID uuid.UUID
}

// Initialise initialises all of the variables in the game
func (game *Game) Initialise() {
	custom := make(map[string]map[string]reflect.Value)
	custom["deltadex/gameplay"] = make(map[string]reflect.Value)
	custom["deltadex/gameplay"]["Ability"] = reflect.ValueOf(Ability{})
	custom["deltadex/gameplay"]["Monster"] = reflect.ValueOf(Monster{})
	custom["deltadex/gameplay"]["Card"] = reflect.ValueOf(Card{})
	custom["deltadex/gameplay"]["Player"] = reflect.ValueOf((*Player)(nil))
	custom["deltadex/gameplay"]["Game"] = reflect.ValueOf((*Game)(nil))
	custom["deltadex/gameplay/events"] = make(map[string]reflect.Value)
	custom["deltadex/gameplay/events"]["Event"] = reflect.ValueOf(events.Event{})
	custom["deltadex/gameplay/events"]["EventID"] = reflect.ValueOf(events.MonsterDieEvent)
	events.LoadScripts(custom)

	game.GameID = uuid.New()

	game.PlayerOne.Hand = []Card{}
	game.PlayerTwo.Hand = []Card{}
	game.PlayerOne.Deck = []Card{}
	game.PlayerTwo.Deck = []Card{}
	game.PlayerOne.Monsters = make([]Monster, 3)
	game.PlayerTwo.Monsters = make([]Monster, 3)

	game.PlayerOne.Energy = 10
	game.PlayerTwo.Energy = 10
	game.PlayerOne.MaxHealth = 50
	game.PlayerTwo.MaxHealth = 50
	game.PlayerOne.Health = game.PlayerOne.MaxHealth
	game.PlayerTwo.Health = game.PlayerTwo.MaxHealth

	game.GameTurn = 0
	game.GameStatus = 0
	game.PlayerTurn = 0
}

// Reset resets the game
func (game *Game) Reset() {
	game.Initialise()
}

// Start starts the game
func (game *Game) Start() {
	// Reset the game before starting
	game.Initialise()
	log.WithFields(log.Fields{
		"game_id": game.GameID,
	}).Info("Starting game, both players connected")

	// Set the game to the first turn and set status to in-progress
	game.GameTurn = 1
	game.GameStatus = 1

	// Decide the starting player
	rand.Seed(time.Now().UnixNano())
	game.PlayerTurn = rand.Intn(2) + 1
	log.WithFields(log.Fields{
		"starting_player": game.PlayerTurn,
		"player_one":      game.PlayerOne.Username,
		"player_two":      game.PlayerTwo.Username,
	}).Info("Starting player decided")

	// Send packets to the players to inform them who starts
	selfContent := map[string]interface{}{
		"username": game.PlayerOne.Username,
		"id":       game.PlayerOne.ID,
		"energy":   game.PlayerOne.Energy,
		"health":   game.PlayerOne.Health,
		"starting": game.PlayerTurn == game.PlayerOne.ID,
	}
	game.PlayerOne.SendPacket(networking.Packet{PacketID: networking.SelfInitiationInformation, Content: selfContent})
	game.PlayerTwo.SendPacket(networking.Packet{PacketID: networking.OpponentInitiationInformation, Content: selfContent})

	selfContent = map[string]interface{}{
		"username": game.PlayerTwo.Username,
		"id":       game.PlayerTwo.ID,
		"energy":   game.PlayerTwo.Energy,
		"health":   game.PlayerTwo.Health,
		"starting": game.PlayerTurn == game.PlayerTwo.ID,
	}
	game.PlayerTwo.SendPacket(networking.Packet{PacketID: networking.SelfInitiationInformation, Content: selfContent})
	game.PlayerOne.SendPacket(networking.Packet{PacketID: networking.OpponentInitiationInformation, Content: selfContent})

	// Send players packets with their starting hands
	card := Card{ID: 0, Name: "Zombie", Type: 0, Attack: 2, Health: 4, EnergyCost: 2, Ability: Ability{AbilityID: 1, Name: "Zombie", Description: "Monster resurrected at half health upon death", Targeted: false}}
	hand := []Card{card, card, card, card}

	packetContent := map[string]interface{}{
		"hand": hand,
	}
	game.PlayerOne.Hand = hand
	game.PlayerTwo.Hand = hand
	game.PlayerOne.SendPacket(networking.Packet{PacketID: networking.StartingHand, Content: packetContent})
	game.PlayerTwo.SendPacket(networking.Packet{PacketID: networking.StartingHand, Content: packetContent})

	game.PlayerOne.SendPacket(networking.Packet{PacketID: networking.OpponentStartingHand, Content: map[string]interface{}{"hand": len(game.PlayerTwo.Hand)}})
	game.PlayerTwo.SendPacket(networking.Packet{PacketID: networking.OpponentStartingHand, Content: map[string]interface{}{"hand": len(game.PlayerOne.Hand)}})
}

// EndTurn is run when each player ends their turn
func (game *Game) EndTurn(playerID int) {
	var player *Player
	if playerID == 1 {
		player = game.PlayerOne
	} else {
		player = game.PlayerTwo
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
			monster := player.OtherPlayer().Monsters[index]
			player.OtherPlayer().Monsters[index] = Monster{}
			events.PushEvent(events.Event{EventID: events.MonsterDieEvent, EventInfo: map[string]interface{}{"game": &CurGame, "monster": monster, "position": index, "player": player.OtherPlayer()}})
		}
	}
}
