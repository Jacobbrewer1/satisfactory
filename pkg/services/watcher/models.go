package watcher

import (
	"encoding/json"
	"time"
)

type vectorMessage struct {
	Message    json.RawMessage `json:"message"`
	Path       string          `json:"path"`
	SourceType string          `json:"source_type"`
	Timestamp  time.Time       `json:"timestamp"`
}

type dockerInfo struct {
	Command      string `json:"Command"`
	CreatedAt    string `json:"CreatedAt"`
	ID           string `json:"ID"`
	Image        string `json:"Image"`
	Labels       string `json:"Labels"`
	LocalVolumes string `json:"LocalVolumes"`
	Mounts       string `json:"Mounts"`
	Names        string `json:"Names"`
	Networks     string `json:"Networks"`
	Ports        string `json:"Ports"`
	RunningFor   string `json:"RunningFor"`
	Size         string `json:"Size"`
	State        string `json:"State"`
	Status       string `json:"Status"`
}
