package main

import (
	"../initbot"
)

var ()

func main() {
	discord := initbot.StartBot()

	defer discord.Close()

	// Let the bot listen to commands
	<-make(chan struct{})
}
