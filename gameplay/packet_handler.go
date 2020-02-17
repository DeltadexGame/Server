package gameplay

import (
	"deltadex/server/networking"
	"deltadex/services"
	"strconv"

	"github.com/Strum355/log"
)

// PacketHandlers keeps track of packet to their handler(s)
var PacketHandlers map[networking.PacketID]func(*Player, networking.Packet) = make(map[networking.PacketID]func(*Player, networking.Packet), 0)

// RegisterPacketHandlers puts all of the handlers into the map
func RegisterPacketHandlers() {
	PacketHandlers[networking.AuthenticationInformation] = handleAuthInfoPacket
	PacketHandlers[networking.PlayerPlayCard] = handlePlayCardPacket
	PacketHandlers[networking.EndTurn] = handleEndTurn
}

func handleAuthInfoPacket(p *Player, packet networking.Packet) {
	info := packet.Content.(map[string]interface{})
	log.WithFields(log.Fields{
		"username": info["username"],
		"token":    info["token"],
	}).Debug("Player information received")

	service := services.DiscordService{}
	service.SendAlert(info["username"].(string) + " attempted to connect to server.")

	deltadex := services.DeltadexService{}
	if deltadex.Authenticate(info["username"].(string), info["token"].(string)) {
		log.Info("User successfully authenticated with API.")
		p.SendPacket(networking.Packet{PacketID: 1, Content: map[string]string{"response": "success"}})
		p.Username = info["username"].(string)
		p.Authenticated = true

		if CurGame.PlayerOne == nil {
			p.ID = 1
			CurGame.PlayerOne = p
			log.WithFields(log.Fields{
				"uuid": CurGame.PlayerOne.UUID,
			}).Info("Player One connected")
		} else {
			p.ID = 2
			CurGame.PlayerTwo = p
			log.WithFields(log.Fields{
				"uuid": CurGame.PlayerTwo.UUID,
			}).Info("Player Two connected")
			CurGame.Start()
		}
	} else {
		log.Info("User attempted to use invalid token.")
		p.SendPacket(networking.Packet{PacketID: 1, Content: map[string]string{"response": "failure"}})
		p.Disconnect()
	}
}

func handlePlayCardPacket(p *Player, packet networking.Packet) {
	info := packet.Content.(map[string]interface{})
	log.WithFields(log.Fields{
		"from": info["from"],
	}).Debug("Card played")

	handPosition, _ := strconv.Atoi(info["from"].(string))
	place, _ := strconv.Atoi(info["place"].(string))

	if p.ID != CurGame.PlayerTurn {
		log.Debug("Player attepted to play card not on their turn")
		packetContent := map[string]interface{}{"result": false, "from": handPosition, "place": place}
		p.SendPacket(networking.Packet{PacketID: networking.PlayCardResult, Content: packetContent})
		return
	}

	if handPosition > len(p.Hand)-1 {
		log.Debug("PLAYER PLAYED CARD NOT IN HAND")
		packetContent := map[string]interface{}{"result": false, "from": handPosition, "place": place}
		p.SendPacket(networking.Packet{PacketID: networking.PlayCardResult, Content: packetContent})
		return
	}

	card := p.Hand[handPosition]
	result := p.PlayCard(card, place)

	packetContent := map[string]interface{}{"result": result, "from": handPosition, "place": place}
	p.SendPacket(networking.Packet{PacketID: networking.PlayCardResult, Content: packetContent})
}

func handleEndTurn(p *Player, packet networking.Packet) {
	if CurGame.PlayerTurn != p.ID {
		return
	}

	CurGame.EndTurn(p)

	if CurGame.NextTurn {
		CurGame.GameTurn++
		CurGame.NextTurn = false
		CurGame.PlayerTurn = p.OtherPlayer().ID
		return
	}

	CurGame.PlayerTurn = p.OtherPlayer().ID
	CurGame.NextTurn = true
}
