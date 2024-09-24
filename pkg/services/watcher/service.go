package watcher

import (
	"context"

	"github.com/Jacobbrewer1/satisfactory/pkg/alerts"
)

type Service interface {
	Start() error
}

type service struct {
	ctx          context.Context
	alertManager alerts.DiscordManager
	listName     string
	alertsURL    string
}

func NewService(ctx context.Context, alertManager alerts.DiscordManager, listName string) Service {
	return &service{
		ctx:          ctx,
		alertManager: alertManager,
		listName:     listName,
	}
}
