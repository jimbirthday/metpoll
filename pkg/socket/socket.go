package socket

import "net"

func TCPSocket(addr string, passive bool) (int, net.Addr, error) {
	return tcpSocket(addr, passive)
}
