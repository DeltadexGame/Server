package networking

// Packet represents a piece of information sent between client and server
type Packet struct {
	PacketID PacketID
	Content  interface{}
}

type PacketID int

const (
	AuthenticationInformation PacketID = 001
	AuthenticationStatus      PacketID = 002

	GameInitiationInformation PacketID = 101
	StartingHand              PacketID = 102
	OpponentStartingHand      PacketID = 103

	PlayerPlayCard   PacketID = 201
	PlayCardResult   PacketID = 202
	OpponentPlayCard PacketID = 203
	RemainingEnergy  PacketID = 204

	EndTurn PacketID = 301
)
