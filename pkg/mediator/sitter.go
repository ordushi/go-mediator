package mediator

import "github.com/google/uuid"

func (b *Observable[T]) AddSitter(e string, ch chan eventMessage[T]) {
	if b.sitters == nil {
		b.sitters = make(map[string][]chan eventMessage[T])
	}
	if _, ok := b.sitters[e]; ok {
		b.sitters[e] = append(b.sitters[e], ch)
	} else {
		b.sitters[e] = []chan eventMessage[T]{ch}
	}
}

func (b *Observable[T]) RemoveSitter(e string, ch chan eventMessage[T]) {
	if _, ok := b.sitters[e]; ok {
		for i := range b.sitters[e] {
			if b.sitters[e][i] == ch {
				b.sitters[e] = append(b.sitters[e][:i], b.sitters[e][i+1:]...)
				break
			}
		}
	}
}

func (b *Observable[T]) Emit(e string, response T) {
	if _, ok := b.sitters[e]; ok {
		for _, handler := range b.sitters[e] {
			go func(handler chan eventMessage[T]) {
				handler <- newEventWrapper(response, false)
			}(handler)
		}
	}
}
func (b *Observable[T]) EmitWithResponse(e string, response T) T {

	request := newEventWrapper(response, true)
	if _, ok := b.sitters[e]; ok {
		for _, handler := range b.sitters[e] {
			go func(handler chan eventMessage[T]) {
				handler <- request
			}(handler)
		}

	}
	return (<-b.Subscriber(request.CorrelationId.String())).Args
}
func newEventWrapper[T any](val T, withresponse bool) eventMessage[T] {
	return eventMessage[T]{withresponse: false, Args: val, CorrelationId: uuid.New()}
}
