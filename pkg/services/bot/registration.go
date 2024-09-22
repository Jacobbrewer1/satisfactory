package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func (s *service) onBotCreate(_ *discordgo.Session, r *discordgo.Ready) {
	slog.Info("Bot is registered as: " + r.User.String())
}
