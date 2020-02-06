package server

import (
	"deltadex/gameplay"
	"fmt"
	"net"

	"github.com/Strum355/log"
	"github.com/spf13/viper"
)

// Server is the main struct for running the Game Server
type Server struct {
}

// Start allows the Server to start receiving connections
func (s *Server) Start() {
	listener, err := net.Listen("tcp", ":"+fmt.Sprint(viper.GetInt("game.tcp.port")))
	if err != nil {
		log.WithError(err).Error("Could not start listening")
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.WithError(err).Error("Could not accept connection")
			continue
		}

		player := gameplay.Player{Conn: conn}
		go player.Handle()
	}
}
