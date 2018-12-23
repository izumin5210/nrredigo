package nrredigo

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/newrelic/go-agent"
)

// Pool is an interface for representing a pool of Redis connections
type Pool interface {
	GetContext(ctx context.Context) (redis.Conn, error)
}

func Wrap(p Pool, opts ...Option) Pool {
	return &wrappedPool{
		Pool: p,
		cfg:  createConfig(opts),
	}
}

type wrappedPool struct {
	Pool
}

func (p *wrappedPool) GetContext(ctx context.Context) (conn redis.Conn, err error) {
	conn, err = p.Pool.GetContext(ctx)
	if err != nil {
		return
	}

	nrtx := newrelic.FromContext(ctx)
	if nrtx != nil {
		conn = wrapConn(conn, nrtx, cfg)
	}

	return
}
