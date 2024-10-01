package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
}

// NewPool returns a new Pool.
func NewPool(poolOpt PoolOption, connOpts ...ConnectionOption) error {
	if poolOpt == nil {
		return errors.New("no pool option provided")
	}

	poolConn := new(pool)
	poolConn.Pool = new(redis.Pool)
	if len(connOpts) != 0 {
		for _, opt := range connOpts {
			opt(poolConn)
		}
	}

	poolOpt(poolConn)

	return nil
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
			slog.Error("error closing connection", slog.String(loggingKeyError, err.Error()))
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
