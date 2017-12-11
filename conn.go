package nrredigo

import (
	"github.com/garyburd/redigo/redis"
	"github.com/newrelic/go-agent"
)

func wrapConn(c redis.Conn, txn newrelic.Transaction) redis.Conn {
	return &wrappedConn{
		Conn: c,
		txn:  txn,
	}
}

type wrappedConn struct {
	redis.Conn
	txn newrelic.Transaction
}

func (c *wrappedConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	// TODO: should start and end a new newrelic datastore segment
	return c.Conn.Do(commandName, args...)
}

func (c *wrappedConn) Send(commandName string, args ...interface{}) error {
	// TODO: should buffer commands and their args
	return c.Conn.Send(commandName, args...)
}

func (c *wrappedConn) Flush() error {
	// TODO: should start and end a new newrelic datastore segment for pipelining operations
	return c.Conn.Flush()
}
