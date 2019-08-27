package messagehandler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler takes a discord message and decides what to do with it
func CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.Bot {
		// Do nothing because the bot is talking
		return
	}

	channelID := message.ChannelID
	channel, _ := session.Channel(channelID)
	server, _ := session.Guild(message.GuildID)
	session.ChannelMessageSend(channelID, "OwO")

	fmt.Printf("(%s > %s) %s: %+v\n", server.Name, channel.Name, message.Author, message.Content)
}
