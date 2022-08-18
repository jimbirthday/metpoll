package metpoll

import "net"

type MSocket struct {
}

func TCPSocket(addr string, passive bool) (int, net.Addr, error) {
	return tcpSocket(addr, passive)
}
