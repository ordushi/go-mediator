package mediator

func (b *Observable[T]) AddSitter(e string, ch chan T) {
	if b.sitters == nil {
		b.sitters = make(map[string][]chan T)
	}
	if _, ok := b.sitters[e]; ok {
		b.sitters[e] = append(b.sitters[e], ch)
	} else {
		b.sitters[e] = []chan T{ch}
	}
}

func (b *Observable[T]) RemoveSitter(e string, ch chan T) {
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
			go func(handler chan T) {
				handler <- response
			}(handler)
		}
	}
}
