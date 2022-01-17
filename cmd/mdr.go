package main

import (
	"fmt"

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
	b := mediator.New[B, string]()
	_ = b.NewMediator("test", test2)
	mtr := a.NewMediator("test", test)
	for i := 0; i < 1000; i++ {
		fmt.Println(
			mtr.Mediate(A{LastName: fmt.Sprint(i), FirstName: fmt.Sprint(i)}))
	}

	fmt.Scanln()

}
func test(tt *mediator.MediatePayload[A, string]) {
	//fmt.Printf(" %s from  - %s", tt.Payload.FirstName, "test1")
	tt.Response = tt.Payload.FirstName
}

func test2(tt *mediator.MediatePayload[B, string]) {
	//	fmt.Printf(" %s from  - %s", tt.Payload.FirstName, "test2")
	tt.Response = "hi2?"
}
