package messagehandler

import (
	"fmt"
	"strings"

	"../accountstats"
	"../config"

	"github.com/bwmarrin/discordgo"
)

var ()

// Validation for well formed bot command
func isValidCommand(message *discordgo.MessageCreate) bool {
	author := message.Author
	inputMessage := message.Content
	inputMessageFields := strings.Fields(inputMessage)

	if author.Bot || len(inputMessageFields) < 2 || !strings.HasPrefix(inputMessage, config.Cfg.BotPrefix) {
		return false
	}
	return true
}

// CommandHandler takes a discord message and decides what to do with it
func CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if isValidCommand(message) == false {
		return
	}

	var outputMessage string
	inputMessage := message.Content[len(config.Cfg.BotPrefix):]
	inputMessageFields := strings.Fields(inputMessage)
	action := inputMessageFields[0]

	channelID := message.ChannelID
	channel, _ := session.Channel(channelID)
	server, _ := session.Guild(message.GuildID)

	switch {
	case action == "sr":
		btags := inputMessageFields[1:]
		emb := accountstats.GetEmbeddedStats(btags)
		session.ChannelMessageSendEmbed(channelID, emb)
	}

	//temp
	// outputMessage = action

	if outputMessage != "" {
		session.ChannelMessageSend(channelID, outputMessage)
	}

	fmt.Printf("(%s > %s) %s: %+v\n", server.Name, channel.Name, message.Author, inputMessage)
}
