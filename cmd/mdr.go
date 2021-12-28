package main

import (
	"fmt"
	"keyless-cache/go-mediator/pkg/mediator"
	"time"
)

type A struct {
	LastName string
}
type B struct {
	Name string
}

func main() {

	//	y := B{s: ""}
	a := mediator.New[A, string]()
	//b := mediator.New[B]()
	mtr := a.NewMediator("a", test)
	go mtr.Listener()
	go mtr.Mediate(A{LastName: "A"})
	time.Sleep(3 * time.Second)

}
func test(tt A) string {
	fmt.Println("111")
	return "yay"
}
