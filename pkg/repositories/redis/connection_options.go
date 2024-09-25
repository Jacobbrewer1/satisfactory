package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

type ConnectionOption func(r *redis.Pool)

func WithMaxIdle(maxIdle int, def ...int) ConnectionOption {
	return func(r *redis.Pool) {
		r.MaxIdle = maxIdle
	}
}

func WithMaxActive(maxActive int) ConnectionOption {
	return func(r *redis.Pool) {
		r.MaxActive = maxActive
	}
}

func WithIdleTimeout(idleTimeout int) ConnectionOption {
	return func(r *redis.Pool) {
		r.IdleTimeout = time.Duration(idleTimeout) * time.Second
	}
}

func WithDial(dial func() (redis.Conn, error)) ConnectionOption {
	return func(r *redis.Pool) {
		r.Dial = dial
	}
}

func FromViper(v *viper.Viper) []ConnectionOption {
	return []ConnectionOption{
		WithMaxIdle(v.GetInt("redis.max_idle")),
		WithMaxActive(v.GetInt("redis.max_active")),
		WithIdleTimeout(v.GetInt("redis.idle_timeout_secs")),
		WithDial(func() (redis.Conn, error) {
			return redis.Dial(
				NetworkTCP,
				v.GetString("redis.address"),
				redis.DialDatabase(v.GetInt("redis.db")),
				redis.DialUsername(v.GetString("redis.username")),
				redis.DialPassword(v.GetString("redis.password")),
			)
		}),
	}
}
