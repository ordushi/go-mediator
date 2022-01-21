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
	//	b := mediator.New[B, string]()
	//	_ = b.NewMediator("test", test2)
	mtr := a.NewMediator("test", test)
	_ = a.NewMediator("test", test2)
	_ = a.NewMediator("test", test2)
	_ = a.NewMediator("test", test2)
	_ = a.NewMediator("test", test2)
	_ = a.NewMediator("test", test2)
	_ = a.NewMediator("test", test2)
	_ = a.NewMediator("test", test2)

	for i := 0; i < 10000; i++ {
		//	go func(mtr mediator.Mediator[A, string], i int) {
		fmt.Println(
			mtr.Mediate(A{LastName: fmt.Sprint(i), FirstName: fmt.Sprint(i)}))
		//}(mtr, i)

	}

	fmt.Scanln()

}
func test(tt *mediator.MediatePayload[A, string]) {
	//fmt.Printf(" %s from  - %s", tt.Payload.FirstName, "test1")
	//fmt.Print("ack: " + tt.Payload.FirstName + " - ")
	tt.Response = tt.Payload.FirstName
}

func test2(tt *mediator.MediatePayload[A, string]) {
	//	fmt.Printf(" %s from  - %s", tt.Payload.FirstName, "test2")
	//	fmt.Print("ack2: " + tt.Payload.FirstName + "  ")
	//	fmt.Print("ack: " + tt.Payload.FirstName + " - ")

	//tt.Response = tt.Payload.FirstName

}
