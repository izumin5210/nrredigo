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

func (c *wrappedConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	if c.txn != nil {
		seg := c.createSegment()
		seg.Operation = commandName
		seg.ParameterizedQuery = formatCommand(commandName, args)
		defer seg.End()
	}
	return c.Conn.Do(commandName, args...)
}

func (c *wrappedConn) Send(commandName string, args ...interface{}) error {
	if c.txn != nil {
		seg := c.createSegment()
		seg.Operation = commandName
		seg.ParameterizedQuery = formatCommand(commandName, args)
		defer seg.End()
	}
	return c.Conn.Send(commandName, args...)
}

func (c *wrappedConn) Flush() error {
	if c.txn != nil {
		seg := c.createSegment()
		seg.Operation = "flush"
		defer seg.End()
	}
	return c.Conn.Flush()
}

func (c *wrappedConn) Receive() (interface{}, error) {
	if c.txn != nil {
		seg := c.createSegment()
		seg.Operation = "receive"
		defer seg.End()
	}
	return c.Conn.Receive()
}

func (c *wrappedConn) createSegment() newrelic.DatastoreSegment {
	return newrelic.DatastoreSegment{
		StartTime: newrelic.StartSegmentNow(c.txn),
		Product:   newrelic.DatastoreRedis,
		// TODO
		// Host:         c.host,
		// PortPathOrID: c.id,
		// DatabaseName: c.databaseName,
	}
}
