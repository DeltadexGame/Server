package networking

// Packet represents a piece of information sent between client and server
type Packet struct {
	PacketID int
	Content  interface{}
}
