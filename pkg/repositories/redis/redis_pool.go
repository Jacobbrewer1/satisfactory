package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	"github.com/gomodule/redigo/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// conn is the global redis connection pool.
	conn Pool

	// ErrRedisNotInitialised is returned when the redis connection pool is not initialised.
	ErrRedisNotInitialised = errors.New("redis connection pool not initialised")
)

// Latency is the duration of Redis queries.
var Latency = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "redis_latency",
		Help: "Duration of Redis queries",
	},
	[]string{"command"},
)

type Pool interface {
	// Do will send a command to the server and returns the received reply on a connection from the pool.
	Do(command string, args ...any) (reply any, err error)

	// DoCtx will send a command to the server with a context and returns the received reply on a connection from the pool.
	DoCtx(ctx context.Context, command string, args ...any) (reply any, err error)

	// Conn returns a redis connection from the pool.
	Conn() redis.Conn
}

// pool represents a redis connection pool.
type pool struct {
	*redis.Pool

	l *slog.Logger
}

// NewPool returns a new Pool.
func NewPool(host string, db int, username, password string) {
	l := slog.With(slog.String(logging.KeyDal, "redis"))

	conn = &pool{
		Pool: &redis.Pool{
			MaxIdle:     3,                 // maximum number of idle connections in the pool (default is 3)
			MaxActive:   0,                 // unlimited connections to the redis server (default is 10 connections)
			IdleTimeout: 240 * time.Second, // 4 minutes idle timeout to match the redis server config (default is 300 seconds)
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial(
					"tcp",
					host,
					redis.DialDatabase(db),
					redis.DialUsername(username),
					redis.DialPassword(password),
				)
				if err != nil {
					return nil, err
				}

				return c, nil
			},
		},
		l: l,
	}
}

// Do will send a command to the server and returns the received reply on a connection from the pool.
func (p *pool) Do(command string, args ...any) (reply any, err error) {
	return p.DoCtx(context.Background(), command, args...)
}

// DoCtx will send a command to the server with a context and returns the received reply on a connection from the pool.
func (p *pool) DoCtx(ctx context.Context, command string, args ...any) (reply any, err error) {
	// The context cannot be nil for the redis pool.
	if ctx == nil {
		ctx = context.Background()
	}

	t := prometheus.NewTimer(Latency.WithLabelValues(command))
	defer t.ObserveDuration()

	c, err := p.GetContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting connection from pool: %w", err)
	}

	defer func(c redis.Conn) {
		if err := c.Close(); err != nil {
			p.l.Error("error closing connection", slog.String(logging.KeyError, err.Error()))
		}
	}(c)

	return c.Do(command, args...)
}

// Conn returns a redis connection from the pool.
func (p *pool) Conn() redis.Conn {
	return p.Pool.Get()
}

func Do(command string, args ...any) (reply any, err error) {
	if conn == nil {
		return nil, ErrRedisNotInitialised
	}
	return DoCtx(context.Background(), command, args...)
}

func DoCtx(ctx context.Context, command string, args ...any) (reply any, err error) {
	if conn == nil {
		return nil, ErrRedisNotInitialised
	}
	return conn.DoCtx(ctx, command, args...)
}

func Conn() redis.Conn {
	if conn == nil {
		return nil
	}
	return conn.Conn()
}
