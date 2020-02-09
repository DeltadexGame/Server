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
	if deltadex.Authenticate(info["token"].(string)) {
		log.Info("User successfully authenticated with API.")
		p.SendPacket(networking.Packet{PacketID: 1, Content: map[string]string{"response": "success"}})
		p.Username = info["username"].(string)
		p.Authenticated = true

		if PlayerOne == nil {
			p.ID = 1
			PlayerOne = p
			log.WithFields(log.Fields{
				"uuid": PlayerOne.UUID,
			}).Info("Player One connected")
		} else {
			p.ID = 2
			PlayerTwo = p
			log.WithFields(log.Fields{
				"uuid": PlayerTwo.UUID,
			}).Info("Player Two connected")
			Start()
		}
	} else {
		log.Info("User attempted to use invalid token.")
		p.SendPacket(networking.Packet{PacketID: 1, Content: map[string]string{"response": "failure"}})
		p.Disconnect()
	}
}

func handlePlayCardPacket(p *Player, packet networking.Packet) {
	if p.ID != PlayerTurn {
		log.Debug("Player attepted to play card not on their turn")
		return
	}

	info := packet.Content.(map[string]interface{})
	log.WithFields(log.Fields{
		"from": info["from"],
	}).Debug("Card played")

	handPosition, _ := strconv.Atoi(info["from"].(string))
	card := p.Hand[handPosition]

	place, _ := strconv.Atoi(info["place"].(string))
	result := p.PlayCard(card, place)

	packetContent := map[string]interface{}{"result": result, "from": handPosition, "place": place}
	p.SendPacket(networking.Packet{PacketID: networking.PlayCardResult, Content: packetContent})
}
