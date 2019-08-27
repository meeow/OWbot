package main

import (
	"../initbot"
)

var ()

func main() {
	//discord, user := initbot.StartBot()

	discord := initbot.StartBot()

	defer discord.Close()

	// Let the bot listen to commands
	<-make(chan struct{})
}
