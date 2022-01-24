package mediator

import (
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
	time        time.Time
	args        T
	name        string
	subscribers *hashmap.HashMap //map[string][]*chan eventMessage[T, K]
}
type eventMessage[T Input, K Output] struct {
	withresponse  bool
	CorrelationId uuid.UUID
	Args          T
	response      K
}

func New[T Input, K Output]() Observable[T, K] {
	return Observable[T, K]{time: time.Now()}
}

func (obs *Observable[T, K]) Subscribe(action string) *chan eventMessage[T, K] {
	ch := make(chan eventMessage[T, K], 0)
	obs.addSubscriber(action, &ch)
	return &ch

}
func returnZero[T any](s ...T) T {
	var zero T
	return zero
}

func (b *Observable[T, K]) addSubscriber(e string, ch *chan eventMessage[T, K]) {
	if b.subscribers == nil {
		//	b.sitters = make(map[string][]*chan eventMessage[T, K])
		b.subscribers = &hashmap.HashMap{}
	}
	if sitter, ok := b.subscribers.Get(e); ok {
		castedSitter := sitter.(*[]*chan eventMessage[T, K])
		*castedSitter = append(*castedSitter, ch)

		return
	} else {
		b.subscribers.Set(e, &[]*chan eventMessage[T, K]{ch})
	}
}

func (b *Observable[T, K]) removeSubscriber(e string, ch *chan eventMessage[T, K]) {
	if sitter, ok := b.subscribers.Get(e); ok {
		castedSitter := *sitter.(*[]*chan eventMessage[T, K])
		for i := range castedSitter {
			if castedSitter[i] == ch {
				castedSitter = append(castedSitter[:i], castedSitter[i+1:]...)
				break
			}
		}
	}
}
func (b *Observable[T, K]) RemoveRSitter(correlationId, e string, ch *chan eventMessage[T, K]) {
	//defer close(*ch)
	b.removeSubscriber(e, ch)
	//b.RemoveSitter(correlationId, ch)
	b.subscribers.Del(correlationId)
	//delete(b.sitters, correlationId)

}

func (b *Observable[T, K]) Emit(e string, request T) {
	if sitter, ok := b.subscribers.Get(e); ok {
		castedSitter := *sitter.(*[]*chan eventMessage[T, K])

		for _, handler := range castedSitter {
			go func(handler chan eventMessage[T, K]) {
				handler <- newEventWrapper[T, K](request, false)
			}(*handler)
		}
	}
}
func (b *Observable[T, K]) EmitResponse(e string, request K) {
	if sitter, ok := b.subscribers.Get(e); ok {
		castedSitter := *sitter.(*[]*chan eventMessage[T, K])
		for _, handler := range castedSitter {
			go func(handler chan eventMessage[T, K]) {
				handler <- eventMessage[T, K]{withresponse: false, response: request}
				defer close(handler)
				defer b.removeSubscriber(e, &handler)
			}(*handler)
		}
	}
}

func (b *Observable[T, K]) EmitWithResponse(e string, requestWrp eventMessage[T, K]) eventMessage[T, K] {

	//requestWrp := newEventWrapper[T, K](request, true)
	if sitter, ok := b.subscribers.Get(e); ok {
		castedSitter := *sitter.(*[]*chan eventMessage[T, K])
		for _, handler := range castedSitter {
			go func(handler chan eventMessage[T, K]) {
				handler <- requestWrp

			}(*handler)
		}

	}
	return requestWrp
}
func newEventWrapper[T Input, K Output](val T, withresponse bool) eventMessage[T, K] {
	return eventMessage[T, K]{withresponse: withresponse, Args: val, CorrelationId: uuid.New()}
}
