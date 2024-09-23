package bot

import "github.com/bwmarrin/discordgo"

const (
	serverInfoCmdID = "server-info"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        serverInfoCmdID,
			Description: "Server Info",
			Type:        discordgo.ChatApplicationCommand,
		},
	}
)
