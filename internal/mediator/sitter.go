package mediator

import "github.com/google/uuid"

func (b *Observable[T, K]) AddSitter(e string, ch chan eventMessage[T, K]) {
	if b.sitters == nil {
		b.sitters = make(map[string][]chan eventMessage[T, K])
	}
	if _, ok := b.sitters[e]; ok {
		b.sitters[e] = append(b.sitters[e], ch)
	} else {
		b.sitters[e] = []chan eventMessage[T, K]{ch}
	}
}

func (b *Observable[T, K]) RemoveSitter(e string, ch chan eventMessage[T, K]) {
	if _, ok := b.sitters[e]; ok {
		for i := range b.sitters[e] {
			if b.sitters[e][i] == ch {
				b.sitters[e] = append(b.sitters[e][:i], b.sitters[e][i+1:]...)
				break
			}
		}
	}
}

func (b *Observable[T, K]) Emit(e string, request T) {
	if _, ok := b.sitters[e]; ok {
		for _, handler := range b.sitters[e] {
			go func(handler chan eventMessage[T, K]) {
				handler <- newEventWrapper[T, K](request, false)
			}(handler)
		}
	}
}
func (b *Observable[T, K]) Response(e string, request K) {
	if _, ok := b.sitters[e]; ok {
		for _, handler := range b.sitters[e] {
			go func(handler chan eventMessage[T, K]) {
				handler <- eventMessage[T, K]{withresponse: false, response: request}
			}(handler)
		}
	}
}

func (b *Observable[T, K]) EmitWithResponse(e string, request T) eventMessage[T, K] {

	requestWrp := newEventWrapper[T, K](request, true)
	if _, ok := b.sitters[e]; ok {
		for _, handler := range b.sitters[e] {
			go func(handler chan eventMessage[T, K]) {
				handler <- requestWrp
			}(handler)
		}

	}
	return requestWrp
}
func newEventWrapper[T Input, K Output](val T, withresponse bool) eventMessage[T, K] {
	return eventMessage[T, K]{withresponse: withresponse, Args: val, CorrelationId: uuid.New()}
}
