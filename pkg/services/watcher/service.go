package watcher

import "context"

type Service interface {
	Start() error
}

type service struct {
	ctx context.Context

	listName string
}

func NewService(ctx context.Context, listName string) Service {
	return &service{
		ctx:      ctx,
		listName: listName,
	}
}
