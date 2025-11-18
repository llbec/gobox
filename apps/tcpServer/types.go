package tcpserver

import (
	"fmt"
	"net"
	"time"
)

// 三元组结构（内部使用）
type triple struct {
	ListenPort int
	PeerIP     string
	PeerPort   int
}

func (t triple) key() string {
	return fmt.Sprintf("%d-%s-%d", t.ListenPort, t.PeerIP, t.PeerPort)
}

type Reservation struct {
	Triple    triple
	Created   time.Time
	LastError error
	Status    string // "ready", "waiting", "error", "disposed"
}

type ConnWrapper struct {
	net.Conn
	Triple triple
	ID     string
	Time   time.Time
}

type EventType int

const (
	EventConnected EventType = iota
	EventDisconnected
)

type ConnEvent struct {
	Type   EventType
	Triple triple
	Conn   *ConnWrapper
	Reason error
}

type Config struct {
	AutoRestartFailed bool
	AutoRestartPeriod time.Duration
}
