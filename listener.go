package metpoll

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/jimbirthday/metpoll/pkg/socket"
	"golang.org/x/sys/unix"
)

type listener struct {
	once    sync.Once
	fd      int
	addr    net.Addr
	address string
}

func (ln *listener) normalize() (err error) {
	ln.fd, ln.addr, err = socket.TCPSocket(ln.address, true)
	return
}

func (ln *listener) close() {
	ln.once.Do(
		func() {
			if ln.fd > 0 {
				fmt.Println(os.NewSyscallError("close", unix.Close(ln.fd)))
			}
		},
	)
}

func initListener(network, addr string) (l *listener, err error) {
	l = &listener{address: addr}
	err = l.normalize()
	return l, err
}
