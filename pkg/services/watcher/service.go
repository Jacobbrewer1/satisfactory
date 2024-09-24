package watcher

import (
	"context"

	"github.com/Jacobbrewer1/satisfactory/pkg/alerts"
)

type Service interface {
	Start() error
}

type service struct {
	ctx                   context.Context
	alertManager          alerts.DiscordManager
	serverInfoListName    string
	serverDetailsListName string
}

func NewService(ctx context.Context, alertManager alerts.DiscordManager, serverInfoListName, serverDetailsListName string) Service {
	return &service{
		ctx:                   ctx,
		alertManager:          alertManager,
		serverInfoListName:    serverInfoListName,
		serverDetailsListName: serverDetailsListName,
	}
}
