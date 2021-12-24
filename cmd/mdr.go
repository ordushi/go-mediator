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
	a := mediator.NewObservable[A]()
	b := mediator.NewObservable[B]()

	go func() {
		for {
			msg := <-a.Subscriber("a")
			fmt.Println("1: " + msg.Args.LastName)
			a.Emit(msg.CorrelationId.String(), A{LastName: "Nami"})

		}
	}()
	time.Sleep(1 * time.Second)

	go func() {
		for {
			msg := <-b.Subscriber("a")
			fmt.Println("2: " + msg.Args.Name)
		}
	}()
	time.Sleep(1 * time.Second)

	fmt.Println(a.EmitWithResponse("a", A{LastName: "Dushi"}))
	b.Emit("a", B{Name: "Or"})
	time.Sleep(3 * time.Second)

}
