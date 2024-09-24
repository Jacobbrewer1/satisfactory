package watcher

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"regexp"

	"github.com/Jacobbrewer1/satisfactory/pkg/repositories/redis"
	redisgo "github.com/gomodule/redigo/redis"
)

func (s *service) processMessage(msg []byte) error {
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
	got, err := redisgo.StringMap(redis.Conn.DoCtx(s.ctx, "HGETALL", "docker_info"))
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

	if _, err := redis.Conn.DoCtx(s.ctx, "HMSET", redisgo.Args{}.Add("docker_info").AddFlat(values)...); err != nil {
		return fmt.Errorf("store docker info: %w", err)
	}

	return nil
}
