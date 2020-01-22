package services

import (
	"bytes"
	"net/http"

	"github.com/Strum355/log"
	"github.com/spf13/viper"
)

// DiscordService used to send alerts through discord
type DiscordService struct {
}

// SendAlert sends an alert through the discord service
func (*DiscordService) SendAlert(message string) {
	if !viper.GetBool("game.production") {
		log.Debug("Discord alert avoided due to non-production")
		return
	}
	url := "https://discordapp.com/api/webhooks/667388689357209610/CoE9bePsCzjlYzoYwLDxo0Whyt6vcdKpv8Js-mbeIsjeHSTNzxwoIyxcQ-p5smC0-Oj6"
	body := []byte(`
	{
		"username": "Game Server",
		"embeds": [{
			"author": {
				"name": "Game Server"
			},
			"title": "Alert",
			"description": "` + message + `"
		}]
	}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.WithError(err).Error("Could not create discord request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.WithError(err).Error("Could not send discord request")
		return
	}
}
