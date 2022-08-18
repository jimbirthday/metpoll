package metpoll

import (
	"syscall"
)

type Poller struct {
	conn map[int]*Connect
	sfd  int
	efd  int
}

func NewPoller() *Poller {
	return &Poller{
		conn: make(map[int]*Connect),
	}
}

func (p *Poller) GetConn(fd int) *Connect {
	return p.conn[fd]
}

func (p *Poller) AddConn(conn *Connect) {
	p.conn[conn.fd] = conn
}

func (p *Poller) DelConn(fd int) {
	delete(p.conn, fd)
}

func (p *Poller) CloseConn(fd int) error {
	conn := p.GetConn(fd)
	if conn == nil {
		return nil
	}
	if err := p.EpollRemoveEvent(fd); err != nil {
		return err
	}
	conn.Close()
	p.DelConn(fd)
	return nil
}

func (p *Poller) EpollRemoveEvent(fd int) error {
	return syscall.EpollCtl(p.efd, syscall.EPOLL_CTL_DEL, fd, nil)
}

func (p *Poller) EpollAddEvent(fd int) error {
	return syscall.EpollCtl(p.efd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(fd),
		Pad:    0,
	})
}

//处理epoll
func (p *Poller) HandlerEpoll() error {
	events := make([]syscall.EpollEvent, 8)
	//在死循环中处理epoll
	for {
		//msec -1,会一直阻塞,直到有事件可以处理才会返回, n 事件个数
		n, err := syscall.EpollWait(p.efd, events, -1)
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			//先在map中是否有这个链接
			conn := p.GetConn(int(events[i].Fd))
			if conn == nil { //没有这个链接,忽略
				continue
			}
			if events[i].Events&syscall.EPOLLHUP == syscall.EPOLLHUP || events[i].Events&syscall.EPOLLERR == syscall.EPOLLERR {
				//断开||出错
				if err := p.CloseConn(int(events[i].Fd)); err != nil {
					return err
				}
			} else if events[i].Events == syscall.EPOLLIN {
				//可读事件
				conn.Read()
			}
		}
	}
}
