package tcpserver

import "time"

type TCPServerAPI interface {
	AddReservationAndListen(listenPort int, remoteIP string, remotePort int) error
	CancelReservation(listenPort int, remoteIP string, remotePort int)
	StartAutoRetry(interval time.Duration)
	StopAutoRetry()
}
