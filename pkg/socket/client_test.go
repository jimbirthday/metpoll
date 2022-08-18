package socket

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestClinet(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("hello"))
	data := make([]byte, 100)
	n, err := conn.Read(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data[:n]))

	time.Sleep(100 * time.Millisecond)
	conn.Close()
}
