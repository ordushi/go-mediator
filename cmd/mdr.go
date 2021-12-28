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
	a := mediator.New[A, (string)]()
	//b := mediator.New[B]()
	mtr := a.NewMediator("test", test)
	mtr2 := a.NewMediator("test", test2)
	go mtr2.Listener()
	time.Sleep(1 * time.Second)

	go mtr.Listener()
	time.Sleep(1 * time.Second)

	fmt.Println(mtr.Mediate(A{LastName: "A"}))

}
func test(tt A) string {
	fmt.Println("111")
	return "yay1"
}
func test2(tt A) string {
	fmt.Println("222")
	return "yay2"
}
