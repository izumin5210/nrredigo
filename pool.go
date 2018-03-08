package nrredigo

import (
	"github.com/garyburd/redigo/redis"
	"github.com/newrelic/go-agent"
)

// Pool is an interface for representing a pool of Redis connections
type Pool interface {
	Get() redis.Conn
}

func WrapPool(p Pool, txn newrelic.Transaction) Pool {
	return &wrappedPool{
		Pool: p,
		txn:  txn,
	}
}

type wrappedPool struct {
	Pool
	txn newrelic.Transaction
}

func (p *wrappedPool) Get() redis.Conn {
	return wrapConn(p.Pool.Get(), p.txn)
}
