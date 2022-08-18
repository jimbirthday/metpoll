package main

import (
	"fmt"

	"github.com/jimbirthday/metpoll"
)

func main() {
	e := metpoll.NewEngine(2)
	err := e.Listen("127.0.0.1:8000")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	e.Run()
}
