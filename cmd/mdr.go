package main

import (
	"fmt"
	"time"

	"ezkv.io/go-mediator/internal/mediator"
)

type A struct {
	FirstName string
	LastName  string
}
type B struct {
	Name string
}

func main() {

	//	y := B{s: ""}
	a := mediator.New[A, string]()
	//b := mediator.New[B]()
	mtr := a.NewMediator("test", test)
	time.Sleep(1 * time.Second)
	//mtr2 := a.NewMediator("test", test2)
	//_ = mtr2
	// go mtr2.Listener()

	// go mtr.Listener()
	time.Sleep(1 * time.Second)
	i := 1

	fmt.Println(
		mtr.Mediate(A{LastName: fmt.Sprint(i), FirstName: fmt.Sprint(i)}))
	//time.Sleep(5 * time.Second)

}
func test(tt *mediator.MediatePayload[A, string]) {
	//fmt.Printf(" %s from  - %s", tt.Payload.FirstName, "test1")
	tt.Response = "hi?"
}

func test2(tt *mediator.MediatePayload[A, string]) {
	//	fmt.Printf(" %s from  - %s", tt.Payload.FirstName, "test2")
	//tt.Response = "hi2?"
}
