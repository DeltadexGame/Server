package config

import "github.com/spf13/viper"

func loadDefaults() {
	viper.SetDefault("game.production", false)
	viper.SetDefault("game.tcp.port", 8080)

}
