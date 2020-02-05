package main

import (
	"weeping-wasp/config"
	"weeping-wasp/gameplay"
	"weeping-wasp/server"
	"weeping-wasp/services"

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
