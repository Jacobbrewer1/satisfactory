package goredis

type PoolOption func(p Pool)

func WithDefaultPool() PoolOption {
	return func(r Pool) {
		conn = r
	}
}

func WithInitializedPool(pool Pool) PoolOption {
	return func(r Pool) {
		conn = pool
	}
}
