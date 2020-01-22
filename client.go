package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	"weeping-wasp/server/server/networking"

	"github.com/Strum355/log"
)

func main() {
	log.InitSimpleLogger(&log.Config{})

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	packet := networking.Packet{PacketID: 0, Content: map[string]string{"username": "oisin", "token": "abcdefg"}}
	err = json.NewEncoder(conn).Encode(packet)
	if err != nil {
		log.WithError(err).Error("connection broke")
		os.Exit(1)
	}
	log.Info("Sent info")

	err = json.NewDecoder(conn).Decode(&packet)
	if err != nil {
		log.WithError(err).Error("connection broke")
		os.Exit(1)
	}
	if packet.PacketID != 1 {
		fmt.Println("packet ID not correct, got ", packet)
	}
	log.Info("Received ack")

	for {
		packet := networking.Packet{PacketID: 2, Content: map[string]string{"username": "oisin", "token": "abcdefg"}}
		err = json.NewEncoder(conn).Encode(packet)
		if err != nil {
			log.WithError(err).Error("connection broke")
			os.Exit(1)
		}
		log.Info("Sent")

		go func() {
			err = json.NewDecoder(conn).Decode(&packet)
			if err != nil {
				log.WithError(err).Error("connection broke")
				os.Exit(1)
			}
			if packet.PacketID != 1 {
				fmt.Println("packet ID not correct, got ", packet)
			}
			log.Info("Received")
		}()
		time.Sleep(5 * time.Second)
	}
}
