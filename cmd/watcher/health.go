package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Jacobbrewer1/goredis/redis"
	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	"github.com/alexliesenfeld/health"
)

func healthHandler() http.Handler {
	checker := health.NewChecker(
		// Disable caching of the results of the checks.
		health.WithCacheDuration(0),
		health.WithDisabledCache(),

		// Set a timeout of 10 seconds for the entire health check.
		health.WithTimeout(10*time.Second),

		// Monitor the health of the database.
		health.WithCheck(health.Check{
			Name: "redis",
			Check: func(ctx context.Context) error {
				_, err := redis.DoCtx(ctx, "PING")
				return err
			},
			Timeout:            3 * time.Second,
			MaxTimeInError:     0,
			MaxContiguousFails: 0,
			StatusListener: func(ctx context.Context, name string, state health.CheckState) {
				logListener(name, state)
			},
			Interceptors:         nil,
			DisablePanicRecovery: false,
		}),
	)

	return health.NewHandler(checker)
}

func logListener(name string, state health.CheckState) {
	switch state.Status {
	case health.StatusUp:
		// The check is healthy.
		slog.Info(fmt.Sprintf("%s is healthy", name))
	case health.StatusDown:
		// The check is unhealthy.
		slog.Error(fmt.Sprintf("%s is unhealthy", name), slog.String(logging.KeyError, state.Result.Error()))
	case health.StatusUnknown:
		// The check is in an unknown state.
		slog.Warn(fmt.Sprintf("%s is in an unknown state", name), slog.String(logging.KeyError, state.Result.Error()))
	default:
		// The check is in an unexpected state.
		slog.Warn(fmt.Sprintf("%s is in an unexpected state", name), slog.String(logging.KeyError, state.Result.Error()))
	}
}
