package metpoll

import (
	"fmt"
	"syscall"
)

type Connect struct {
	fd int
	sd syscall.Sockaddr
}

func (c *Connect) Close() {
	err := syscall.Close(c.fd)
	if err != nil {
		fmt.Printf("fd%d close error:%s\n", c.fd, err.Error())
	}
}

//读取数据
func (c *Connect) Read() {
	data := make([]byte, 1024)

	//通过系统调用,读取数据,n是读到的长度
	n, err := syscall.Read(c.fd, data)
	if n == 0 {
		return
	}
	if err != nil {
		fmt.Printf("fd %d read error:%s\n", c.fd, err.Error())
	} else {
		fmt.Printf("%d say: %s \n", c.fd, data[:n])
		c.Write([]byte(fmt.Sprintf("hello %d", c.fd)))
	}
}

func (c *Connect) Write(data []byte) {
	_, err := syscall.Write(c.fd, data)
	if err != nil {
		fmt.Printf("fd %d write error:%s\n", c.fd, err.Error())
	}
}
