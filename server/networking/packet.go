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

	SelfInitiationInformation     PacketID = 101
	OpponentInitiationInformation PacketID = 102
	StartingHand                  PacketID = 103
	OpponentStartingHand          PacketID = 104

	PlayerPlayCard   PacketID = 201
	PlayCardResult   PacketID = 202
	OpponentPlayCard PacketID = 203
	RemainingEnergy  PacketID = 204

	EndTurn PacketID = 301
)
