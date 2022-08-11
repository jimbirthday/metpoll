package socket

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestServer(t *testing.T) {
	_, _, err := TCPSocket("127.0.0.1:8080", true)
	if err != nil {
		fmt.Println("TCPSocket:" + err.Error())
		return
	}
	fd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		fmt.Println("EpollCreate1:" + err.Error())
		return
	}
	//没有设置 EFD_NONBLOCK 就是阻塞得，设置的话就是非租塞得返回-1
	efd, err := unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EPOLL_CLOEXEC)
	if err != nil {
		fmt.Println("EpollCreate1:" + err.Error())
		return
	}

	unix.EpollCtl(fd, unix.EPOLL_CTL_ADD, efd, &unix.EpollEvent{Fd: int32(efd), Events: unix.EPOLLOUT})

	go func() {
		msec := -1
		ls := make([]unix.EpollEvent, 128)
		for {
			n, err := unix.EpollWait(fd, ls, msec)
			if n == 0 || (n < 0 && err == unix.EINTR) {
				msec = -1
				runtime.Gosched()
				continue
			} else if err != nil {
				fmt.Println("EpollWait:" + err.Error())
				return
			}

			msec = 0

			for i := 0; i < n; i++ {
				ev := &ls[i]
				if fd := int(ev.Fd); fd != efd {
					fmt.Println("EpollWait efd")
				} else { // poller is awakened to run tasks in queues.
					efdBuf := make([]byte, 8)
					_, _ = unix.Read(efd, efdBuf)
					s := string(efdBuf)
					if s != "\x00\x00\x00\x00\x00\x00\x00\x00" {
						fmt.Println(s)
					}
				}
			}
		}
	}()
	time.Sleep(time.Minute * 60)
}

// func TestTCP(t *testing.T) {
// 	listen, err := net.Listen("tcp", "127.0.0.1:8080")
// 	if err != nil {
// 		fmt.Println("Listen() failed, err: ", err)
// 		return
// 	}
// 	for {
// 		conn, err := listen.Accept() // 监听客户端的连接请求
// 		if err != nil {
// 			fmt.Println("Accept() failed, err: ", err)
// 			continue
// 		}
// 		go process(conn) // 启动一个goroutine来处理客户端的连接请求
// 	}
// }

// func process(conn net.Conn) {
// 	defer conn.Close() // 关闭连接
// 	for {
// 		reader := bufio.NewReader(conn)
// 		var buf [128]byte
// 		n, err := reader.Read(buf[:]) // 读取数据
// 		if err != nil {
// 			fmt.Println("read from client failed, err: ", err)
// 			break
// 		}
// 		recvStr := string(buf[:n])
// 		fmt.Println("收到Client端发来的数据：", recvStr)
// 		conn.Write([]byte(recvStr)) // 发送数据
// 	}
// }
