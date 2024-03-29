package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	configFilePath = "./config.toml"
	// Cfg contains all fields from the configFile.
	Cfg = ReadConfig()
)

type config struct {
	BotPrefix     string
	BotStatus     string
	TokenFilePath string
	StatsKeys     []string
}

// ReadConfig reads info from config file
func ReadConfig() config {
	_, err := os.Stat(configFilePath)
	if err != nil {
		log.Fatal("Config file is missing: ", configFilePath)
	}

	var config config
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		log.Fatal(err)
	}

	log.Print("Loaded from config file: ", config)

	return config
}
