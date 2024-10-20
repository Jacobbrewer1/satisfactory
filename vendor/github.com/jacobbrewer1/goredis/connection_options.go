package goredis

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

func WithAddress(address string) ConnectionOption {
	return func(r *pool) {
		r.addr = address
	}
}

func WithNetwork(network string) ConnectionOption {
	return func(r *pool) {
		r.network = network
	}
}

func WithDialOpts(dialOpts ...redis.DialOption) ConnectionOption {
	return func(r *pool) {
		r.dialOpts = dialOpts
	}
}

func FromViper(v *viper.Viper) []ConnectionOption {
	return []ConnectionOption{
		WithMaxIdle(v.GetInt(viperMaxIdle)),
		WithMaxActive(v.GetInt(viperMaxActive)),
		WithIdleTimeout(v.GetInt(viperIdleTimeout)),
		WithAddress(v.GetString(viperAddress)),
		WithNetwork(NetworkTCP),
		WithDialOpts(
			redis.DialDatabase(v.GetInt(viperDatabase)),
			redis.DialUsername(v.GetString(viperUsername)),
			redis.DialPassword(v.GetString(viperPassword)),
		),
	}
}
