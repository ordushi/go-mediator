package mediator

import "github.com/google/uuid"

func (b *Observable[T]) AddSitter(e string, ch chan eventMessage[interface{}]) {
	if b.sitters == nil {
		b.sitters = make(map[string][]chan eventMessage[interface{}])
	}
	if _, ok := b.sitters[e]; ok {
		b.sitters[e] = append(b.sitters[e], ch)
	} else {
		b.sitters[e] = []chan eventMessage[interface{}]{ch}
	}
}

func (b *Observable[T]) RemoveSitter(e string, ch chan eventMessage[interface{}]) {
	if _, ok := b.sitters[e]; ok {
		for i := range b.sitters[e] {
			if b.sitters[e][i] == ch {
				b.sitters[e] = append(b.sitters[e][:i], b.sitters[e][i+1:]...)
				break
			}
		}
	}
}

func (b *Observable[T]) Emit(e string, response interface{}) {
	if _, ok := b.sitters[e]; ok {
		for _, handler := range b.sitters[e] {
			go func(handler chan eventMessage[interface{}]) {
				handler <- newEventWrapper(response, false)
			}(handler)
		}
	}
}

func (b *Observable[T]) EmitWithResponse(e string, response interface{}) eventMessage[interface{}] {

	request := newEventWrapper(response, true)
	if _, ok := b.sitters[e]; ok {
		for _, handler := range b.sitters[e] {
			go func(handler chan eventMessage[interface{}]) {
				handler <- request
			}(handler)
		}

	}
	return request
}
func newEventWrapper[T any](val T, withresponse bool) eventMessage[T] {
	return eventMessage[T]{withresponse: false, Args: val, CorrelationId: uuid.New()}
}
