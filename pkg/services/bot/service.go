package bot

import "github.com/bwmarrin/discordgo"

type Service interface {
	// Start starts the bot
	Start() error

	// Stop stops the bot
	Stop() error
}

type service struct {
	token string
	s     *discordgo.Session
}

func NewService(token string) Service {
	return &service{
		token: token,
	}
}
