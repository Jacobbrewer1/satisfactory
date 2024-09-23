package bot

import (
	"fmt"
	"log/slog"

	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	"github.com/bwmarrin/discordgo"
)

func (s *service) Start() error {
	session, err := discordgo.New("Bot " + s.token)
	if err != nil {
		return fmt.Errorf("failed to create discord session: %w", err)
	}
	s.s = session

	slog.Debug("Registering handlers")
	s.registerHandlers()
	slog.Debug("Handlers registered")

	err = s.s.Open()
	if err != nil {
		return fmt.Errorf("failed to open discord session: %w", err)
	}

	s.removeAllCommands()

	slog.Debug("Registering commands")
	s.shutdownFunc, err = s.registerCommands()
	if err != nil {
		return fmt.Errorf("failed to register commands: %w", err)
	}
	slog.Debug("Commands registered")

	return nil
}

func (s *service) removeAllCommands() {
	slog.Debug("Removing all commands")
	cmds, err := s.s.ApplicationCommands(s.s.State.User.ID, "")
	if err != nil {
		slog.Error("Failed to get commands", slog.String(logging.KeyError, err.Error()))
		return
	}

	for _, v := range cmds {
		err = s.s.ApplicationCommandDelete(s.s.State.User.ID, "", v.ID)
		if err != nil {
			slog.Error("Failed to delete command", slog.String("command", v.Name), slog.String(logging.KeyError, err.Error()))
		}
	}
	slog.Debug("All commands removed")
}

func (s *service) registerHandlers() {
	s.s.AddHandler(s.onBotCreate)
	s.s.AddHandler(s.onInteractionCreate)

	s.interactionHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		serverInfoCmdID: s.onServerInfo,
	}
}

func (s *service) registerCommands() (func(), error) {
	// Register commands here
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	removeCommands := func() {
		s.removeRegisteredCommands(registeredCommands)
	}

	for i, v := range commands {
		cmd, err := s.s.ApplicationCommandCreate(s.s.State.User.ID, "", v) // GuildID is empty because we are creating global commands
		if err != nil {
			removeCommands()
			return nil, fmt.Errorf("cannot create '%v' command: %w", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	return removeCommands, nil
}

func (s *service) removeRegisteredCommands(registeredCommands []*discordgo.ApplicationCommand) {
	// Register commands here
	for _, v := range registeredCommands {
		err := s.s.ApplicationCommandDelete(s.s.State.User.ID, "", v.ID) // GuildID is empty because we are deleting global commands
		if err != nil {
			slog.Error("Failed to delete command", slog.String("command", v.Name), slog.String(logging.KeyError, err.Error()))
		}
	}
}

func (s *service) Stop() error {
	s.shutdownFunc()
	return s.s.Close()
}
