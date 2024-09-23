package watcher

import (
	"fmt"
	"log/slog"

	"github.com/Jacobbrewer1/satisfactory/pkg/repositories/redis"
	redisgo "github.com/gomodule/redigo/redis"
)

func (s *service) Start() error {
	for {
		select {
		case <-s.ctx.Done():
			return fmt.Errorf("context done: %w", s.ctx.Err())
		default:
			got, err := redisgo.ByteSlices(redis.Conn.DoCtx(s.ctx, "BLPOP", s.listName, 0))
			if err != nil {
				return fmt.Errorf("redis lpop: %w", err)
			} else if got == nil {
				slog.Debug("No message to process")
				continue
			}

			if err := s.processMessage(got[1]); err != nil {
				return fmt.Errorf("process message: %w", err)
			}
		}
	}
}
