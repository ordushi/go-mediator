package mediator

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type Observable[T any] struct {
	time    time.Time
	args    T
	name    string
	sitters map[string][]chan eventMessage[interface{}]
}
type eventMessage[T any] struct {
	withresponse  bool
	CorrelationId uuid.UUID
	Args          T
	response      interface{}
}

func New[T any]() Observable[T] {
	return Observable[T]{time: time.Now()}
}
func (obs *Observable[T]) Handle(ctx context.Context) {
	log.Printf("deleting: %+v\n", obs)

}
func (obs *Observable[T]) Subscriber(action string) chan eventMessage[interface{}] {
	ch := make(chan eventMessage[interface{}])
	obs.AddSitter(action, ch)
	return ch

}
