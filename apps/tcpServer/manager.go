package tcpserver

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type Manager struct {
	cfg          Config
	mu           sync.RWMutex
	reservations map[string]*Reservation
	listeners    map[int]*listenerHolder
	events       chan ConnEvent
	stop         chan struct{}
	wg           sync.WaitGroup
}

func NewManager(cfg Config) *Manager {
	return &Manager{
		cfg:          cfg,
		reservations: map[string]*Reservation{},
		listeners:    map[int]*listenerHolder{},
		events:       make(chan ConnEvent, 128),
		stop:         make(chan struct{}),
	}
}

func (m *Manager) Events() <-chan ConnEvent {
	return m.events
}

func (m *Manager) Start() {
	if m.cfg.AutoRestartFailed {
		m.wg.Add(1)
		go m.autoRestartLoop()
	}
}

func (m *Manager) Stop() {
	close(m.stop)
	m.wg.Wait()

	m.mu.Lock()
	for _, l := range m.listeners {
		close(l.stop)
		l.ln.Close()
	}
	m.mu.Unlock()
}

func (m *Manager) AddService(listenPort int, peerIP string, peerPort int) (*Reservation, error) {
	tr := triple{listenPort, peerIP, peerPort}
	key := tr.key()

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.reservations[key]; exists {
		return m.reservations[key], nil
	}

	for _, r := range m.reservations {
		if r.Triple.ListenPort == listenPort {
			if peerIP == "*" || r.Triple.PeerIP == "*" {
				return nil, errors.New("peer conflict: port=0 means ANY peers")
			}
			if r.Triple.PeerIP == peerIP && (peerPort == 0 || r.Triple.PeerPort == 0 || r.Triple.PeerPort == peerPort) {
				return nil, errors.New("peer already reserved")
			}
		}
	}

	res := &Reservation{
		Triple:  tr,
		Created: time.Now(),
		Status:  "waiting",
	}
	m.reservations[key] = res

	err := m.ensureListener(listenPort)
	if err != nil {
		res.Status = "error"
		res.LastError = err
		return res, err
	}

	res.Status = "ready"
	return res, nil
}

func (m *Manager) ensureListener(port int) error {
	if holder, ok := m.listeners[port]; ok && holder.running {
		return nil
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	holder := &listenerHolder{
		port:    port,
		ln:      ln,
		running: true,
		stop:    make(chan struct{}),
	}

	m.listeners[port] = holder
	go runListener(m, holder)

	return nil
}

func (m *Manager) RemoveService(listenPort int, peerIP string, peerPort int) {
	tr := triple{listenPort, peerIP, peerPort}
	key := tr.key()

	m.mu.Lock()
	delete(m.reservations, key)
	m.mu.Unlock()

	m.cleanupListener(listenPort)
}

func (m *Manager) cleanupListener(port int) {
	for _, r := range m.reservations {
		if r.Triple.ListenPort == port {
			return
		}
	}

	if holder, ok := m.listeners[port]; ok {
		close(holder.stop)
		holder.ln.Close()
		delete(m.listeners, port)
	}
}

func (m *Manager) autoRestartLoop() {
	ticker := time.NewTicker(m.cfg.AutoRestartPeriod)
	defer ticker.Stop()
	defer m.wg.Done()

	for {
		select {
		case <-m.stop:
			return
		case <-ticker.C:
			m.RetryFailedOnce()
		}
	}
}

func (m *Manager) RetryFailedOnce() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, r := range m.reservations {
		if r.Status != "error" {
			continue
		}
		err := m.ensureListener(r.Triple.ListenPort)
		if err == nil {
			r.Status = "ready"
			r.LastError = nil
		}
	}
}
