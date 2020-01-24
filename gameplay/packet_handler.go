package gameplay

import (
	"weeping-wasp/server/server/networking"
	"weeping-wasp/server/services"

	"github.com/Strum355/log"
)

// PacketHandlers keeps track of packet to their handler(s)
var PacketHandlers map[int]func(*Player, networking.Packet) = make(map[int]func(*Player, networking.Packet), 0)

// RegisterPacketHandlers puts all of the handlers into the map
func RegisterPacketHandlers() {
	PacketHandlers[0] = handleAuthInfoPacket
	PacketHandlers[2] = handleKeepAlivePacket
}

func handleAuthInfoPacket(p *Player, packet networking.Packet) {
	info := packet.Content.(map[string]interface{})
	log.WithFields(log.Fields{
		"username": info["username"],
		"token":    info["token"],
	}).Info("Player information received")

	wasp := services.WaspService{}
	if wasp.Authenticate(info["token"].(string)) {
		log.Info("User successfully authenticated with API.")
		p.SendPacket(networking.Packet{PacketID: 1, Content: map[string]string{"response": "success"}})

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

func handleKeepAlivePacket(p *Player, packet networking.Packet) {
	log.WithFields(log.Fields{
		"player": p.ID,
	}).Info("Keep alive received")
	p.SendPacket(networking.Packet{PacketID: 2, Content: "success"})
}
