package socket

import (
	"errors"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

var listenerBacklogMaxSize = maxListenerBacklog()

func tcpSocket(addr string, passive bool) (int, net.Addr, error) {
	sa, tAddr, err := GetTcpSocket(addr)
	if err != nil {
		return 0, nil, err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, unix.IPPROTO_TCP)
	if err != nil {
		return 0, nil, err
	}

	defer func() {
		// ignore EINPROGRESS for non-blocking socket connect, should be processed by caller
		if err != nil {
			if err, ok := err.(*os.SyscallError); ok && err.Err == unix.EINPROGRESS {
				return
			}
			_ = unix.Close(fd)
		}
	}()

	if passive {
		if err = os.NewSyscallError("bind", unix.Bind(fd, sa)); err != nil {
			return 0, nil, err
		}
		if err = os.NewSyscallError("listen", unix.Listen(fd, listenerBacklogMaxSize)); err != nil {
			return 0, nil, err
		}
	} else {
		if err = os.NewSyscallError("connect", unix.Connect(fd, sa)); err != nil {
			return 0, nil, err
		}
	}
	return fd, tAddr, err
}

func GetTcpSocket(addr string) (sa unix.Sockaddr, tcpAddr *net.TCPAddr, err error) {
	tAddr, err := net.ResolveTCPAddr("", addr)
	if err != nil {
		return
	}
	if tAddr.IP == nil {
		return nil, nil, errors.New("tcp addr is empty")
	}
	sa4 := &unix.SockaddrInet4{
		Port: tAddr.Port,
	}
	copy(sa4.Addr[:], tAddr.IP)

	return sa4, tAddr, nil
}
