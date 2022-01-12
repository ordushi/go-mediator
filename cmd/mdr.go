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
	//b := mediator.New[B]()
	mtr := a.NewMediator("test", test)
	i := 1
	mtr.Stop()
	fmt.Println(
		mtr.Mediate(A{LastName: fmt.Sprint(i), FirstName: fmt.Sprint(i)}))
	fmt.Println(
		mtr.Mediate(A{LastName: fmt.Sprint(i), FirstName: fmt.Sprint(i)}))

	//time.Sleep(5 * time.Second)

}
func test(tt *mediator.MediatePayload[A, string]) {
	//fmt.Printf(" %s from  - %s", tt.Payload.FirstName, "test1")
	//	tt.Response = "hi?"
}

func test2(tt *mediator.MediatePayload[A, string]) {
	//	fmt.Printf(" %s from  - %s", tt.Payload.FirstName, "test2")
	//tt.Response = "hi2?"
}
