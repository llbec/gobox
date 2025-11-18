package tcpserver

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type listenerHolder struct {
	port    int
	ln      net.Listener
	running bool
	stop    chan struct{}
	wg      sync.WaitGroup
}

func runListener(mgr *Manager, holder *listenerHolder) {
	holder.wg.Add(1)
	defer holder.wg.Done()

	for {
		conn, err := holder.ln.Accept()
		if err != nil {
			select {
			case <-holder.stop:
				return
			default:
			}
			holder.running = false
			return
		}

		remoteIP, remotePort, _ := net.SplitHostPort(conn.RemoteAddr().String())

		// 匹配三元组
		tr := mgr.matchReservation(remoteIP, remotePort, holder.port)
		if tr == nil {
			log.Println("No reservation matched for incoming connection:", remoteIP)
			conn.Close()
			continue
		}

		cw := &ConnWrapper{
			Conn:   conn,
			Triple: *tr,
			Time:   time.Now(),
			ID:     fmt.Sprintf("%d", time.Now().UnixNano()),
		}

		mgr.events <- ConnEvent{
			Type:   EventConnected,
			Triple: *tr,
			Conn:   cw,
		}

		go func() {
			buf := make([]byte, 2048)
			for {
				_, err := conn.Read(buf)
				if err != nil {
					mgr.events <- ConnEvent{
						Type:   EventDisconnected,
						Triple: *tr,
						Reason: err,
					}
					conn.Close()
					return
				}
			}
		}()
	}
}

func (mgr *Manager) matchReservation(remoteIP, remotePort string, listenPort int) *triple {
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()

	for _, r := range mgr.reservations {
		if r.Triple.ListenPort != listenPort {
			continue
		}
		if r.Triple.PeerIP != remoteIP && r.Triple.PeerIP != "*" {
			continue
		}
		if fmt.Sprintf("%d", r.Triple.PeerPort) != remotePort && r.Triple.PeerPort != 0 {
			continue
		}
		return &r.Triple
	}
	return nil
}
