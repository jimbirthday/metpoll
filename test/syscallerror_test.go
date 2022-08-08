package test

import (
	"fmt"
	"os"
	"testing"
)

func TestSyscall(t *testing.T) {
	err := os.NewSyscallError("test", a())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("done")

}
func a() error {

	return nil
}
