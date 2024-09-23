package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func (s *service) onInteractionCreate(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	handler, ok := s.interactionHandlers[i.ApplicationCommandData().Name]
	if !ok {
		slog.Error("No handler found for command", slog.String("command", i.ApplicationCommandData().Name))
		return
	}

	handler(s.s, i)
}
