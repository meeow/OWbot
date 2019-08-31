package main

import (
	"../initapi"
	"../initbot"
)

var ()

func main() {
	initapi.StartAllEndpoints()
	discord := initbot.StartBot()

	defer discord.Close()

	// Let the bot listen to commands
	<-make(chan struct{})
}
