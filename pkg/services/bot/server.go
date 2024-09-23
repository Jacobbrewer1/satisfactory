package bot

import (
	"context"
	"log/slog"
	"time"

	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	"github.com/Jacobbrewer1/satisfactory/pkg/repositories/redis"
	"github.com/Jacobbrewer1/satisfactory/pkg/utils"
	"github.com/bwmarrin/discordgo"
	redisgo "github.com/gomodule/redigo/redis"
)

func (s *service) onServerInfo(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	// Respond to the user with "Just getting the server info"
	err := s.s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("Error responding to server info", slog.String(logging.KeyError, err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the server info from redis
	serverInfo, err := redisgo.StringMap(redis.Conn.DoCtx(ctx, "HGETALL", "docker_info"))
	if err != nil {
		slog.Error("Error getting server info", slog.String(logging.KeyError, err.Error()))
		return
	}

	// Send the server info to the user
	msg := "State: " + serverInfo["State"] + "\n" +
		"RunningFor: " + serverInfo["RunningFor"] + "\n" +
		"Status: " + serverInfo["Status"]

	_, err = s.s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: utils.Ptr(msg),
	})
	if err != nil {
		slog.Error("Error deleting server info", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) onServerCredentials(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	// Respond to the user with "Just getting the server credentials"
	err := s.s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("Error responding to server credentials", slog.String(logging.KeyError, err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the server credentials from redis
	serverCredentials, err := redisgo.StringMap(redis.Conn.DoCtx(ctx, "HGETALL", "server_credentials"))
	if err != nil {
		slog.Error("Error getting server credentials", slog.String(logging.KeyError, err.Error()))
		return
	}

	// Send the server credentials to the user
	msg := "IP: " + serverCredentials["ip"] + "\n" +
		"Port: " + serverCredentials["port"] + "\n" +
		"Password: " + serverCredentials["password"]

	_, err = s.s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: utils.Ptr(msg),
	})
	if err != nil {
		slog.Error("Error deleting server credentials", slog.String(logging.KeyError, err.Error()))
		return
	}
}
