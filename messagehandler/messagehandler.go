package messagehandler

import (
	"fmt"
	"strings"

	"../accountstats"
	"../config"

	"github.com/bwmarrin/discordgo"
)

const (
	minimumFields = 2
)

// Validation for well formed bot command
func isValidCommand(message *discordgo.MessageCreate) bool {
	author := message.Author
	inputMessage := message.Content
	inputMessageFields := strings.Fields(inputMessage)

	if author.Bot ||
		!strings.HasPrefix(inputMessage, config.Cfg.BotPrefix) ||
		len(inputMessageFields) < minimumFields {
		return false
	}
	return true
}

// CommandHandler takes a discord message and decides what to do with it
func CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if isValidCommand(message) == false {
		return // The message does not invoke a bot action
	}

	outputMessage := ""
	var outputEmbeds []*discordgo.MessageEmbed
	inputMessage := message.Content[len(config.Cfg.BotPrefix):]
	inputMessageFields := strings.Fields(inputMessage)
	action := strings.ToLower(inputMessageFields[0])

	channelID := message.ChannelID
	channel, _ := session.Channel(channelID)
	server, _ := session.Guild(message.GuildID)

	// improve logging later
	fmt.Printf("(%s > %s) %s: %+v\n", server.Name, channel.Name, message.Author, inputMessage)

	switch {
	case action == "sr":
		btags := inputMessageFields[1:]
		outputEmbeds = append(outputEmbeds, accountstats.GetAllEmbeddedStats(btags)...)
	}

	if outputMessage != "" {
		session.ChannelMessageSend(channelID, outputMessage)
	}

	for _, emb := range outputEmbeds {
		go session.ChannelMessageSendEmbed(channelID, emb)
	}

}
