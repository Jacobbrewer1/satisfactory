package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

type ConnectionOption func(r *pool)

func WithMaxIdle(maxIdle int) ConnectionOption {
	return func(r *pool) {
		r.MaxIdle = maxIdle
	}
}

func WithMaxActive(maxActive int) ConnectionOption {
	return func(r *pool) {
		r.MaxActive = maxActive
	}
}

func WithIdleTimeout(idleTimeout int) ConnectionOption {
	return func(r *pool) {
		r.IdleTimeout = time.Duration(idleTimeout) * time.Second
	}
}

func WithDial(dial func() (redis.Conn, error)) ConnectionOption {
	return func(r *pool) {
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
