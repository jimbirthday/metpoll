package socket

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	_, _, err := TCPSocket("127.0.0.1:8080", true)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
