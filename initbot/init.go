package initbot

import (
	"fmt"
	"io/ioutil"
	"strings"

	"../config"
	"../messagehandler"
	"github.com/bwmarrin/discordgo"
)

var ()

func getToken(path string) string {
	token, err := ioutil.ReadFile(config.Cfg.TokenFilePath)
	if err != nil {
		fmt.Println(err)
	}
	strToken := strings.TrimSpace(string(token))

	return strToken
}

// StartBot gets the bot token and uses it to start a discord session.
// It also adds desired handlers.
func StartBot() *discordgo.Session {
	discord, err := discordgo.New("Bot " + getToken(config.Cfg.TokenFilePath))

	discord.AddHandler(messagehandler.CommandHandler)

	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, config.Cfg.BotStatus)
		servers := discord.State.Guilds
		fmt.Println("OWbot has started on servers:")
		for _, server := range servers {
			fmt.Println(server.Name, server.ID) // server.Name does not currently work, API may be broken
		}
		fmt.Println(len(servers), "in total")
	})

	err = discord.Open()
	fmt.Println(err)
	return discord
}
