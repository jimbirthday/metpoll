package socket

import (
	"fmt"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	_, _, err := TCPSocket("127.0.0.1:8080", true)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	time.Sleep(time.Second * 10)
}
