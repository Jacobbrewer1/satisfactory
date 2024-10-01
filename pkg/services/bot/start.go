package bot

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/Jacobbrewer1/goredis/redis"
	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	"github.com/bwmarrin/discordgo"
	redisgo "github.com/gomodule/redigo/redis"
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

	go s.handleBotStatus()

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
		serverInfoCmdID:        s.onServerInfo,
		serverCredentialsCmdID: s.onServerCredentials,
		severDetailsCmdID:      s.onServerDetails,
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

func (s *service) handleBotStatus() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	currentNum := -1 // Start at -1 so the bot status is updated on the first iteration

	for range ticker.C {
		num, err := s.getPlayersConnected()
		if err != nil {
			slog.Error("Failed to get players connected", slog.String(logging.KeyError, err.Error()))
			continue
		}

		if num == currentNum {
			continue
		}

		currentNum = num

		if err := s.setBotStatusPlayerCount(num); err != nil {
			slog.Error("Failed to update bot status", slog.String(logging.KeyError, err.Error()))
			continue
		}

		slog.Debug("Bot status updated")
	}
}

func (s *service) setBotStatusPlayerCount(num int) error {
	return s.s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: fmt.Sprintf("%d players", num),
				Type: discordgo.ActivityTypeWatching,
				URL:  "",
			},
		},
	})
}

func (s *service) getPlayersConnected() (int, error) {
	// Get players connected to the server from redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	got, err := redisgo.StringMap(redis.DoCtx(ctx, "HGETALL", "server_details"))
	if err != nil {
		return 0, fmt.Errorf("failed to get players connected: %w", err)
	}

	playersConnected, err := strconv.Atoi(got["NumConnectedPlayers"])
	if err != nil {
		return 0, fmt.Errorf("failed to convert players connected to int: %w", err)
	}

	return playersConnected, nil
}
