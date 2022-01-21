package mediator

import (
	"sync"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
)

type Input interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string | ~struct{} | interface{}
}
type Output interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string | ~struct{} | interface{}
}
type Observable[T Input, K Output] struct {
	time       time.Time
	args       T
	name       string
	sitters    *hashmap.HashMap //map[string][]*chan eventMessage[T, K]
	responsers []chan eventMessage[T, K]
	mutex      sync.RWMutex
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

func (obs *Observable[T, K]) Subscriber(action string) *chan eventMessage[T, K] {
	ch := make(chan eventMessage[T, K], 0)
	obs.AddSitter(action, &ch)
	return &ch

}
func returnZero[T any](s ...T) T {
	var zero T
	return zero
}
