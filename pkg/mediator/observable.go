package mediator

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

var flavio chan string = make(chan string)

type Observable[T any] struct {
	time    time.Time
	args    T
	name    string
	sitters map[string][]chan eventMessage[T]
}
type eventMessage[T any] struct {
	withresponse  bool
	CorrelationId uuid.UUID
	Args          T
}

func (obs *Observable[T]) Handle(ctx context.Context) {
	log.Printf("deleting: %+v\n", obs)

}
func (obs *Observable[T]) Subscriber(action string) chan eventMessage[T] {
	ch := make(chan eventMessage[T])
	obs.AddSitter(action, ch)
	return ch

}
func NewObservable[T any]() Observable[T] {
	return new[T]()

}
