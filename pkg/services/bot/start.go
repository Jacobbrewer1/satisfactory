package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (s *service) Start() error {
	session, err := discordgo.New("Bot " + s.token)
	if err != nil {
		return fmt.Errorf("failed to create discord session: %w", err)
	}

	s.s = session

	err = s.s.Open()
	if err != nil {
		return fmt.Errorf("failed to open discord session: %w", err)
	}

	return nil
}

func (s *service) RegisterHandlers() {
	s.s.AddHandler(s.onBotCreate)
}

func (s *service) Stop() error {
	return s.s.Close()
}
