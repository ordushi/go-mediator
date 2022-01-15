package mediator

import (
	"sync"

	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
)

func (b *Observable[T, K]) AddSitter(e string, ch *chan eventMessage[T, K]) {
	if b.sitters == nil {
		//	b.sitters = make(map[string][]*chan eventMessage[T, K])
		b.sitters = &hashmap.HashMap{}
	}
	if sitter, ok := b.sitters.Get(e); ok {
		castedSitter := sitter.(*[]*chan eventMessage[T, K])
		*castedSitter = append(*castedSitter, ch)
		//b.sitters.Set(e, castedSitter)

		return
		//b.sitters[e] = append(b.sitters[e], ch)
	} else {

		b.sitters.Set(e, &[]*chan eventMessage[T, K]{ch})
		//b.sitters[e] = []*chan eventMessage[T, K]{ch}
	}
}

func (b *Observable[T, K]) RemoveSitter(e string, ch *chan eventMessage[T, K]) {
	if sitter, ok := b.sitters.Get(e); ok {
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
	defer close(*ch)
	b.RemoveSitter(e, ch)
	b.RemoveSitter(correlationId, ch)
	go b.sitters.Del(correlationId)
	//delete(b.sitters, correlationId)

}

func (b *Observable[T, K]) Emit(e string, request T) {
	if sitter, ok := b.sitters.Get(e); ok {
		castedSitter := *sitter.(*[]*chan eventMessage[T, K])

		for _, handler := range castedSitter {
			go func(handler chan eventMessage[T, K]) {
				handler <- newEventWrapper[T, K](request, false)
			}(*handler)
		}
	}
}
func (b *Observable[T, K]) Response(e string, request K) {
	if sitter, ok := b.sitters.Get(e); ok {
		castedSitter := *sitter.(*[]*chan eventMessage[T, K])
		for _, handler := range castedSitter {
			go func(handler chan eventMessage[T, K], mutex *sync.RWMutex) {
				handler <- eventMessage[T, K]{withresponse: false, response: request}
			}(*handler, &b.mutex)
		}
	}
}

func (b *Observable[T, K]) EmitWithResponse(e string, request T) eventMessage[T, K] {

	requestWrp := newEventWrapper[T, K](request, true)
	if sitter, ok := b.sitters.Get(e); ok {
		castedSitter := *sitter.(*[]*chan eventMessage[T, K])

		for _, handler := range castedSitter {
			go func(handler chan eventMessage[T, K], mutex *sync.RWMutex) {
				handler <- requestWrp

			}(*handler, &b.mutex)
		}

	}
	return requestWrp
}
func newEventWrapper[T Input, K Output](val T, withresponse bool) eventMessage[T, K] {
	return eventMessage[T, K]{withresponse: withresponse, Args: val, CorrelationId: uuid.New()}
}
