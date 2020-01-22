package config

import (
	"encoding/json"
	"github.com/Strum355/log"
	"github.com/spf13/viper"
	"strings"
)

// Load loads the env variables into the config
func Load() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	loadDefaults()
	viper.AutomaticEnv()
}

// PrintSettings prints out all the settings in a debug message
func PrintSettings() {
	settings := viper.AllSettings()

	out, _ := json.MarshalIndent(settings, "", "\t")
	log.Debug("config:\n" + string(out))
}
