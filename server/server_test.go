package server

import (
	"bufio"
	"encoding/json"
	"net"
	"strings"
	"testing"
	"time"
	"weeping-wasp/server/server/networking"

	"github.com/Strum355/log"
)

func TestRun(t *testing.T) {
	log.InitSimpleLogger(&log.Config{})

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(reply) != "hello new connection :)" {
		t.Fatal("reply not the right one, got: ", reply)
	}

	packet := networking.Packet{PacketID: 0, Content: map[string]string{"username": "oisin", "token": "abcdefg"}}
	json.NewEncoder(conn).Encode(packet)

	json.NewDecoder(conn).Decode(&packet)
	if packet.PacketID != 1 {
		t.Fatal("packet ID not correct, got ", packet)
	}
}
