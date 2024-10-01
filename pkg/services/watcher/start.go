package watcher

import (
	"context"
	"log/slog"

	"github.com/Jacobbrewer1/goredis/redis"
	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	redisgo "github.com/gomodule/redigo/redis"
)

func (s *service) Start() error {
	go s.watchServerInfo(s.ctx)
	go s.watchServerDetails(s.ctx)

	<-s.ctx.Done()

	return nil
}

func (s *service) watchServerInfo(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("Context done")
			return
		default:
			got, err := redisgo.ByteSlices(redis.DoCtx(ctx, "BLPOP", s.serverInfoListName, 0))
			if err != nil {
				slog.Error("Error getting message from redis info list", slog.String(logging.KeyError, err.Error()))
				continue
			} else if got == nil {
				slog.Debug("No message to process")
				continue
			}

			if err := s.processInfoMessage(got[1]); err != nil {
				slog.Error("Error processing info message", slog.String(logging.KeyError, err.Error()))
			}

			slog.Debug("Info message processed")
		}
	}
}

func (s *service) watchServerDetails(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("Context done")
			return
		default:
			got, err := redisgo.ByteSlices(redis.DoCtx(ctx, "BLPOP", s.serverDetailsListName, 0))
			if err != nil {
				slog.Error("Error getting message from redis details list", slog.String(logging.KeyError, err.Error()))
				continue
			} else if got == nil {
				slog.Debug("No message to process")
				continue
			}

			if err := s.processDetailsMessage(got[1]); err != nil {
				slog.Error("Error processing details message", slog.String(logging.KeyError, err.Error()))
			}

			slog.Debug("Details message processed")
		}
	}
}
