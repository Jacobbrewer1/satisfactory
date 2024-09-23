package bot

import "github.com/bwmarrin/discordgo"

const (
	serverInfoCmdID        = "server-info"
	serverCredentialsCmdID = "server-credentials"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        serverInfoCmdID,
			Type:        discordgo.ChatApplicationCommand,
			Description: "Server Info",
		},
		{
			Name:        serverCredentialsCmdID,
			Type:        discordgo.ChatApplicationCommand,
			Description: "Server Credentials",
		},
	}
)
