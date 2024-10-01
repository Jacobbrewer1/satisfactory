package watcher

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"regexp"
	"strconv"

	"github.com/Jacobbrewer1/goredis/redis"
	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	redisgo "github.com/gomodule/redigo/redis"
)

func (s *service) processInfoMessage(msg []byte) error {
	vecMsg := new(vectorMessage)
	if err := json.Unmarshal(msg, vecMsg); err != nil {
		return fmt.Errorf("unmarshal vector message: %w", err)
	}

	// Regex remove all backslashes unless they are \"
	reg := regexp.MustCompile(`\\([^"])`)
	vecMsg.Message = []byte(reg.ReplaceAllString(string(vecMsg.Message), ""))

	// Replace \"\" with \" to fix the unmarshal
	regSlash := regexp.MustCompile(`\\"\\"`)
	vecMsg.Message = []byte(regSlash.ReplaceAllString(string(vecMsg.Message), "\""))

	// Replace \" with " to fix the unmarshal
	regQuote := regexp.MustCompile(`\\"`)
	vecMsg.Message = []byte(regQuote.ReplaceAllString(string(vecMsg.Message), "\""))

	// Remove all whitespace
	regSpace := regexp.MustCompile(`\s+`)
	vecMsg.Message = []byte(regSpace.ReplaceAllString(string(vecMsg.Message), ""))

	// Remove all newlines
	regNewline := regexp.MustCompile(`\n+`)
	vecMsg.Message = []byte(regNewline.ReplaceAllString(string(vecMsg.Message), ""))

	// Remove all tabs
	regTab := regexp.MustCompile(`\t+`)
	vecMsg.Message = []byte(regTab.ReplaceAllString(string(vecMsg.Message), ""))

	// Remove the first and last character
	vecMsg.Message = vecMsg.Message[1 : len(vecMsg.Message)-1]

	docInfo := new(dockerInfo)
	if err := json.Unmarshal(vecMsg.Message, docInfo); err != nil {
		return fmt.Errorf("unmarshal docker info: %w", err)
	}

	if err := s.handleDockerInfo(*docInfo); err != nil {
		return fmt.Errorf("handle docker info: %w", err)
	}

	return nil
}

// Store and process the message
func (s *service) handleDockerInfo(info dockerInfo) error {
	// Get the current hash map of docker info
	got, err := redisgo.StringMap(redis.DoCtx(s.ctx, "HGETALL", "docker_info"))
	if err != nil {
		return fmt.Errorf("get docker info: %w", err)
	}

	// Compare the new info to the old info
	if got["State"] != info.State {
		slog.Debug("State changed", slog.String("old", got["State"]), slog.String("new", info.State))
		err = s.alertManager.SendDiscordAlert(fmt.Sprintf("Server state changed from `%s` to `%s`", got["State"], info.State))
		if err != nil {
			return fmt.Errorf("send discord alert: %w", err)
		}
	}

	// Store all the info
	v := reflect.ValueOf(info)
	values := make([]any, v.NumField()*2)
	for i := 0; i < v.NumField(); i++ {
		values[i*2] = v.Type().Field(i).Name
		values[i*2+1] = v.Field(i).String()
	}

	if _, err := redis.DoCtx(s.ctx, "HMSET", redisgo.Args{}.Add("docker_info").AddFlat(values)...); err != nil {
		return fmt.Errorf("store docker info: %w", err)
	}

	return nil
}

func (s *service) processDetailsMessage(msg []byte) error {
	vecMsg := new(vectorMessage)
	if err := json.Unmarshal(msg, vecMsg); err != nil {
		return fmt.Errorf("unmarshal vector message: %w", err)
	}

	// Regex remove all backslashes unless they are \"
	reg := regexp.MustCompile(`\\([^"])`)
	vecMsg.Message = []byte(reg.ReplaceAllString(string(vecMsg.Message), ""))

	// Replace \"\" with \" to fix the unmarshal
	regSlash := regexp.MustCompile(`\\"\\"`)
	vecMsg.Message = []byte(regSlash.ReplaceAllString(string(vecMsg.Message), "\""))

	// Replace \" with " to fix the unmarshal
	regQuote := regexp.MustCompile(`\\"`)
	vecMsg.Message = []byte(regQuote.ReplaceAllString(string(vecMsg.Message), "\""))

	// Remove all whitespace
	regSpace := regexp.MustCompile(`\s+`)
	vecMsg.Message = []byte(regSpace.ReplaceAllString(string(vecMsg.Message), ""))

	// Remove all newlines
	regNewline := regexp.MustCompile(`\n+`)
	vecMsg.Message = []byte(regNewline.ReplaceAllString(string(vecMsg.Message), ""))

	// Remove all tabs
	regTab := regexp.MustCompile(`\t+`)
	vecMsg.Message = []byte(regTab.ReplaceAllString(string(vecMsg.Message), ""))

	// Remove the first and last character
	vecMsg.Message = vecMsg.Message[1 : len(vecMsg.Message)-1]

	details := new(serverDetails)
	if err := json.Unmarshal(vecMsg.Message, details); err != nil {
		return fmt.Errorf("unmarshal docker info: %w", err)
	}

	if err := s.handleServerDetails(*details.Data.ServerGameState); err != nil {
		return fmt.Errorf("handle docker info: %w", err)
	}

	return nil
}

func (s *service) handleServerDetails(details ServerGameState) error {
	// Get the current hash map of server details
	got, err := redisgo.StringMap(redis.DoCtx(s.ctx, "HGETALL", "server_details"))
	if err != nil {
		return fmt.Errorf("get server details: %w", err)
	}

	// Compare the new details to the old details
	if got["ActiveSessionName"] != details.ActiveSessionName {
		slog.Debug("Active session name changed", slog.String("old", got["ActiveSessionName"]), slog.String("new", details.ActiveSessionName))
		err = s.alertManager.SendDiscordAlert(fmt.Sprintf("Active session name changed from `%s` to `%s`", got["ActiveSessionName"], details.ActiveSessionName))
		if err != nil {
			return fmt.Errorf("send discord alert: %w", err)
		}
	}

	gameRunning, err := strconv.ParseBool(got["IsGameRunning"])
	if err != nil {
		slog.Error("Error parsing bool", slog.String("value", got["IsGameRunning"]), slog.String(logging.KeyError, err.Error()))
	} else if gameRunning != details.IsGameRunning {
		slog.Debug("Game running changed", slog.Bool("old", gameRunning), slog.Bool("new", details.IsGameRunning))
		err = s.alertManager.SendDiscordAlert(fmt.Sprintf("Game running changed from `%t` to `%t`", gameRunning, details.IsGameRunning))
		if err != nil {
			return fmt.Errorf("send discord alert: %w", err)
		}
	}

	gamePaused, err := strconv.ParseBool(got["IsGamePaused"])
	if err != nil {
		slog.Error("Error parsing bool", slog.String("value", got["IsGamePaused"]), slog.String(logging.KeyError, err.Error()))
	} else if gamePaused != details.IsGamePaused {
		slog.Debug("Game paused changed", slog.Bool("old", gamePaused), slog.Bool("new", details.IsGamePaused))
		err = s.alertManager.SendDiscordAlert(fmt.Sprintf("Game paused changed from `%t` to `%t`", gamePaused, details.IsGamePaused))
		if err != nil {
			return fmt.Errorf("send discord alert: %w", err)
		}
	}

	// Store all the details
	v := reflect.ValueOf(details)
	values := make([]any, v.NumField()*2)
	for i := 0; i < v.NumField(); i++ {
		values[i*2] = v.Type().Field(i).Name

		switch v.Type().Field(i).Type.Kind() {
		case reflect.Int:
			values[i*2+1] = v.Field(i).Int()
		case reflect.Float64:
			values[i*2+1] = v.Field(i).Float()
		case reflect.String:
			values[i*2+1] = v.Field(i).String()
		case reflect.Bool:
			values[i*2+1] = v.Field(i).Bool()
		default:
			return fmt.Errorf("unknown type: %s", v.Type().Field(i).Type.Kind())
		}
	}

	if _, err := redis.DoCtx(s.ctx, "HMSET", redisgo.Args{}.Add("server_details").AddFlat(values)...); err != nil {
		return fmt.Errorf("store server details: %w", err)
	}

	return nil
}
