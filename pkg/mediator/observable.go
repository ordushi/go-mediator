package mediator

import (
	"context"
	"log"
	"time"
)

var flavio chan string = make(chan string)

type Observable[T any] struct {
	time    time.Time
	args    T
	name    string
	sitters map[string][]chan T
}

func (obs *Observable[T]) Handle(ctx context.Context) {
	log.Printf("deleting: %+v\n", obs)

}
func (obs *Observable[T]) Subscriber(action string) chan T {
	ch := make(chan T)
	obs.AddSitter(action, ch)
	return ch

}
func NewObservable[T any]() Observable[T] {
	return new[T]()

}
