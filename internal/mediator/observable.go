package mediator

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Input interface {
	comparable
}
type Output interface {
	comparable
}
type Observable[T Input, K Output] struct {
	time    time.Time
	args    T
	name    string
	sitters map[string][]chan eventMessage[T, K]
	mutex   sync.RWMutex
}
type eventMessage[T Input, K Output] struct {
	withresponse  bool
	CorrelationId uuid.UUID
	Args          T
	response      K
}

func New[T Input, K Output]() Observable[T, K] {
	return Observable[T, K]{time: time.Now(), mutex: sync.RWMutex{}}
}

func (obs *Observable[T, K]) Subscriber(action string) chan eventMessage[T, K] {
	ch := make(chan eventMessage[T, K], 1)
	obs.AddSitter(action, ch)
	return ch

}
func returnZero[T any](s ...T) T {
	var zero T
	return zero
}
