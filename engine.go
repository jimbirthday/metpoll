package metpoll

import (
	"fmt"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

type engine struct {
	pollers []*Poller
	sfd     int
	addr    net.Addr
	lb      LB
}

func NewEngine(ps int) *engine {
	pls := make([]*Poller, ps)
	return &engine{
		pollers: pls,
		lb: &RoundRobinBalance{
			total: 2,
			cur:   1,
		},
	}
}

func (e *engine) Listen(addr string) error {
	fd, ad, err := TCPSocket(addr, true)
	if err != nil {
		return err
	}
	e.sfd = fd
	e.addr = ad

	err = e.createPoller()
	if err != nil {
		return err
	}
	return nil
}

func (e *engine) createPoller() error {
	for i := 0; i < len(e.pollers); i++ {
		efd, err := syscall.EpollCreate1(unix.EPOLL_CLOEXEC)
		if err != nil {
			return err
		}
		cs := make(map[int]*Connect)
		p := &Poller{
			sfd:  e.sfd,
			efd:  efd,
			conn: cs,
		}
		e.pollers[i] = p
		go p.HandlerEpoll()

	}

	fmt.Println("1", e.pollers)
	return nil
}

func (e *engine) Run() {
	for {
		nfd, sd, err := syscall.Accept(e.sfd)
		if err != nil {
			continue
		}
		i := e.lb.Cur()
		fmt.Println("lb ", i)
		p := e.pollers[i]
		p.EpollAddEvent(nfd)
		p.AddConn(&Connect{
			fd: nfd,
			sd: sd,
		})
	}
}
