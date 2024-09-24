package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordManager interface {
	SendDiscordAlert(message string) error
}

type discordPayload struct {
	Content string `json:"content"`
}

type discordManager struct {
	webhookURL string
}

func NewDiscordManager(webhookURL string) DiscordManager {
	return &discordManager{
		webhookURL: webhookURL,
	}
}

func (d *discordManager) SendDiscordAlert(message string) error {
	bdy := &discordPayload{
		Content: message,
	}

	bdyBytes, err := json.Marshal(bdy)
	if err != nil {
		return fmt.Errorf("failed to marshal discord payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, d.webhookURL, bytes.NewBuffer(bdyBytes))
	if err != nil {
		return fmt.Errorf("failed to create discord request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discord request: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord request failed: %s", resp.Status)
	}

	return nil
}
