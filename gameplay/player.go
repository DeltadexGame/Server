package gameplay

import (
	"deltadex/server/networking"
	"encoding/json"
	"math/rand"
	"net"
	"time"

	"github.com/Strum355/log"
	"github.com/google/uuid"
)

// Player holds information about an individual connected to the server
type Player struct {
	ID            int
	Username      string
	Conn          net.Conn
	UUID          uuid.UUID
	connected     bool
	Monsters      []Monster
	Hand          []Card
	Deck          []Card
	Energy        int
	Health        int
	MaxHealth     int
	Authenticated bool
}

// Handle takes an incoming connection and deals with it
func (p *Player) Handle() {
	p.UUID = uuid.New()
	log.WithFields(log.Fields{
		"uuid": p.UUID,
	}).Info("Assigned UUID to incoming connection")

	p.connected = true

	go p.listen()
}

// listen runs a loop while a player is connected, receiving their packets
func (p *Player) listen() {
	for p.connected {
		var packet networking.Packet
		err := json.NewDecoder(p.Conn).Decode(&packet)
		if err != nil {
			p.Conn.Close()
			log.WithFields(log.Fields{
				"player": p.ID,
				"error":  err.Error(),
			}).Info("Player disconnected")
			break
		}

		p.PacketReceived(packet)
	}
}

// PacketReceived is the method called on packet received
func (p *Player) PacketReceived(packet networking.Packet) {
	if packet.PacketID != networking.AuthenticationInformation && !p.Authenticated {
		return
	}

	fun, ok := PacketHandlers[packet.PacketID]
	if !ok {
		log.WithFields(log.Fields{
			"id":      packet.PacketID,
			"content": packet.Content,
		}).Info("Unknown packet received")
		return
	}
	fun(p, packet)
}

// SendPacket sends a packet to the player
func (p *Player) SendPacket(packet networking.Packet) error {
	return json.NewEncoder(p.Conn).Encode(packet)
}

// OtherPlayer returns the Player that is not this player
func (p *Player) OtherPlayer() *Player {
	if p.ID == 1 {
		return CurGame.PlayerTwo
	}
	return CurGame.PlayerOne
}

func (p *Player) GenerateHand(size int) {
	p.shuffleDeck()
	p.Hand = p.Deck[len(p.Deck)-size-1 : len(p.Deck)-1]
	p.Deck = append([]Card{}, p.Deck[:len(p.Deck)-size-1]...)
}

func (p *Player) shuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p.Deck), func(i, j int) { p.Deck[i], p.Deck[j] = p.Deck[j], p.Deck[i] })
}

// Disconnect terminates with the connection with that player
func (p *Player) Disconnect() {
	p.Conn.Close()
	log.WithFields(log.Fields{
		"uuid": p.UUID,
	}).Info("User disconnected")
}

// SpawnMonster creates a monster at the given position
func (p *Player) SpawnMonster(monster Monster, position int) {
	p.Monsters[position] = monster
	p.SendPacket(networking.Packet{PacketID: networking.MonsterSpawn, Content: map[string]interface{}{"position": position, "ownership": true, "monster": monster}})
	p.OtherPlayer().SendPacket(networking.Packet{PacketID: networking.MonsterSpawn, Content: map[string]interface{}{"position": position, "ownership": false, "monster": monster}})
}

// Damage damages the player by a given amount
func (p *Player) Damage(damage int) {
	p.Health -= damage
}

// DrawCard draws a random card from their deck
func (p *Player) DrawCard() {
	// DRAW RANDOM CARD FROM DECK

	rand.Seed(time.Now().UnixNano())
	card := p.Deck[len(p.Deck)-1]
	p.Deck = append([]Card{}, p.Deck[:len(p.Deck)-1]...)
	p.Hand = append(p.Hand, card)

	p.SendPacket(networking.Packet{PacketID: networking.DrawCard, Content: map[string]interface{}{"card": card}})
	p.OtherPlayer().SendPacket(networking.Packet{PacketID: networking.OpponentDrawCard, Content: map[string]interface{}{}})
}

// PlayCard plays the selected card to the selected position
func (p *Player) PlayCard(card Card, position int) bool {
	if card.EnergyCost > p.Energy {
		return false
	}

	if card.Type != 0 && p.Monsters[position] != (Monster{}) {
		return false
	}

	p.Energy -= card.EnergyCost

	packetContent := map[string]interface{}{"card": card, "position": position}
	p.OtherPlayer().SendPacket(networking.Packet{PacketID: networking.OpponentPlayCard, Content: packetContent})

	p.SendPacket(networking.Packet{PacketID: networking.RemainingEnergy, Content: map[string]interface{}{"energy": p.Energy}})

	p.removeFromHand(position)

	if card.Type == 0 {
		p.SpawnMonster(Monster{card.Name, card.Attack, card.Health, card.Health, card.Ability}, position)
	}
	return true
}

func (p *Player) removeFromHand(position int) {
	newHand := []Card{}
	for x, y := range p.Hand {
		if x != position {
			newHand = append(newHand, y)
		}
	}
	p.Hand = newHand
}
