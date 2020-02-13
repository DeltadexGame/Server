package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Strum355/log"
	"github.com/spf13/viper"
)

// DeltadexService is used to talk to the backend API for the game
type DeltadexService struct {
}

// Authenticate verifies the given token is valid
func (w *DeltadexService) Authenticate(username string, token string) bool {
	if !viper.GetBool("game.production") {
		return true
	}
	send := struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}{username, token}

	payload, err := json.Marshal(send)
	if err != nil {
		log.WithError(err).Error("Could not marshal payload")
		return false
	}

	resp, err := http.Post(viper.GetString("game.api")+"/auth", "application/json", strings.NewReader(string(payload)))
	if err != nil {
		log.WithError(err).Error("Error connecting to API")
		return false
	}

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"status_code": resp.StatusCode,
		}).Info("Response not 200.")
		return false
	}
	return true
}
