package nrredigo

import (
	"github.com/garyburd/redigo/redis"
)

// Pool is an interface for representing a pool of Redis connections
type Pool interface {
	Get() redis.Conn
}
