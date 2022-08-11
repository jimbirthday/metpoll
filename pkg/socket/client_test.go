package socket

import (
	"fmt"
	"os"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestClinet(t *testing.T) {
	sfd, _, err := TCPSocket("127.0.0.1:8080", false)
	if err, ok := err.(*os.SyscallError); ok && err.Err == unix.EINPROGRESS {
	} else {
		fmt.Println("TCPSocket:" + err.Error())
		return
	}
	time.Sleep(time.Second * 10)
	fmt.Println("fd:", sfd)
	bys := []byte("hello world")
	_, err = unix.Write(sfd, bys)
	if err != nil {
		fmt.Println("Write:" + err.Error())
		return
	}
	unix.Close(sfd)
}
