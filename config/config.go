package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	configFile = "./config.toml"
	Cfg        = ReadConfig()
)

type Config struct {
	BotPrefix string
	BotStatus string
}

// ReadConfig reads info from config file
func ReadConfig() Config {
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Config file is missing: ", configFile)
	}

	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	log.Print("Loaded configs: ", config)

	return config
}
