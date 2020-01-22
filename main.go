package main

import (
	"weeping-wasp/server/config"
	"weeping-wasp/server/gameplay"
	"weeping-wasp/server/server"
	"weeping-wasp/server/services"

	"github.com/Strum355/log"
)

func main() {
	// Setup logger
	log.InitSimpleLogger(&log.Config{})
	log.Info("Server starting...")

	// Load config
	config.Load()
	config.PrintSettings()

	// Setup discord service
	discord := services.DiscordService{}
	discord.SendAlert("Starting Game Server")

	// Register packet handlers
	gameplay.RegisterPacketHandlers()

	// Start the server
	server := server.Server{}
	server.Start()
}
