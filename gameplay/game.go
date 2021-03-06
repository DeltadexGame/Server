package gameplay

import (
	"deltadex/gameplay/events"
	"deltadex/server/networking"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Strum355/log"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
)

var (
	CurGame Game         = Game{}
	Cards   map[int]Card = make(map[int]Card)
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
	if viper.GetBool("game.download") {
		downloadCards()
	}
	loadCards()
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

func RegisterHandlers() {
	events.RegisterHandler(events.MonsterDamageEvent, HandleHeavyAbility)
	events.RegisterHandler(events.MonsterDieEvent, HandleZombieAbility)
}

// Start starts the game
func (game *Game) Start() {
	// Reset the game before starting
	game.Initialise()
	log.WithFields(log.Fields{
		"game_id": game.GameID,
	}).Info("Starting game, both players connected")

	RegisterHandlers()

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
	for i := 0; i < 30; i++ {
		rand.Seed(time.Now().UnixNano())
		game.PlayerOne.Deck = append(game.PlayerOne.Deck, Cards[rand.Intn(len(Cards))+1])
	}
	for i := 0; i < 30; i++ {
		rand.Seed(time.Now().UnixNano())
		game.PlayerTwo.Deck = append(game.PlayerTwo.Deck, Cards[rand.Intn(len(Cards))+1])
	}

	game.PlayerOne.GenerateHand(4)
	game.PlayerOne.SendPacket(networking.Packet{PacketID: networking.StartingHand, Content: map[string]interface{}{"hand": game.PlayerOne.Hand}})
	game.PlayerTwo.GenerateHand(4)
	game.PlayerTwo.SendPacket(networking.Packet{PacketID: networking.StartingHand, Content: map[string]interface{}{"hand": game.PlayerTwo.Hand}})

	game.PlayerOne.SendPacket(networking.Packet{PacketID: networking.OpponentStartingHand, Content: map[string]interface{}{"hand": len(game.PlayerTwo.Hand)}})
	game.PlayerTwo.SendPacket(networking.Packet{PacketID: networking.OpponentStartingHand, Content: map[string]interface{}{"hand": len(game.PlayerOne.Hand)}})

	if game.PlayerTurn == 1 {
		game.PlayerOne.SendPacket(networking.Packet{PacketID: networking.StartTurn, Content: map[string]interface{}{"self": "true"}})
		game.PlayerTwo.SendPacket(networking.Packet{PacketID: networking.StartTurn, Content: map[string]interface{}{"self": "false"}})
	} else {
		game.PlayerOne.SendPacket(networking.Packet{PacketID: networking.StartTurn, Content: map[string]interface{}{"self": "false"}})
		game.PlayerTwo.SendPacket(networking.Packet{PacketID: networking.StartTurn, Content: map[string]interface{}{"self": "true"}})
	}
}

// EndTurn is run when each player ends their turn
func (game *Game) EndTurn(player *Player) {
	player.OtherPlayer().SendPacket(networking.Packet{PacketID: networking.EndTurn, Content: map[string]interface{}{"self": false}})
	player.SendPacket(networking.Packet{PacketID: networking.EndTurn, Content: map[string]interface{}{"self": true}})
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

		attacked := player.OtherPlayer().Monsters[index]
		damageEvent := events.Event{EventID: events.MonsterDamageEvent, EventInfo: map[string]interface{}{"game": &CurGame, "monster": attacked, "attacker": monster, "position": index, "player": player.OtherPlayer(), "cancelled": false, "damage": monster.Attack}}
		damageEvent = events.PushEvent(damageEvent)
		if !damageEvent.EventInfo["cancelled"].(bool) {
			player.OtherPlayer().Monsters[index].Damage(damageEvent.EventInfo["damage"].(int))
		}
		died := false
		if player.OtherPlayer().Monsters[index].Health <= 0 {
			event := events.Event{EventID: events.MonsterDieEvent, EventInfo: map[string]interface{}{"game": &CurGame, "monster": attacked, "position": index, "player": player.OtherPlayer(), "cancelled": false}}
			event = events.PushEvent(event)
			if !event.EventInfo["cancelled"].(bool) {
				player.OtherPlayer().Monsters[index] = Monster{}
				died = true
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
	}

	player.OtherPlayer().DrawCard()

	player.OtherPlayer().SendPacket(networking.Packet{PacketID: networking.StartTurn, Content: map[string]interface{}{"self": "true"}})
	player.SendPacket(networking.Packet{PacketID: networking.StartTurn, Content: map[string]interface{}{"self": "false"}})
}

func downloadCards() {
	os.RemoveAll(".cache/")
	_, err := git.PlainClone(".cache/Cards", false, &git.CloneOptions{
		URL:      "https://github.com/DeltadexGame/Cards",
		Progress: os.Stdout,
	})
	if err != nil {
		log.WithError(err).Error("Could not download cards. Aborting!!")
		os.Exit(1)
	}
}

func loadCards() {
	if _, err := os.Stat(".cache/"); err != nil {
		if os.IsNotExist(err) {
			downloadCards()
		}
	}
	files, err := ioutil.ReadDir(".cache/Cards/cards")
	if err != nil {
		log.WithError(err).Error("Could not load cards")
		os.Exit(1)
		return
	}

	for _, file := range files {
		f, _ := os.Open(".cache/Cards/cards/" + file.Name())
		card := Card{}
		json.NewDecoder(f).Decode(&card)
		result, err := strconv.Atoi(strings.Replace(file.Name(), ".json", "", 1))
		if err != nil {
			log.WithError(err).Error("Could not get ID of card.")
			continue
		}
		Cards[result] = card
	}

	log.WithFields(log.Fields{
		"cards": len(Cards),
	}).Info("Loaded cards")
}
