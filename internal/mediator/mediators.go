package mediator

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type Mediator[T Input, K Output] struct {
	action     func(*MediatePayload[T, K])
	observable *Observable[T, K]
	actionName string
}
type MediatePayload[T Input, K Output] struct {
	Payload  T
	Response K
}
type IMediator interface {
	Mediate()
	Forfit()
}

func (obs *Observable[T, K]) NewMediator(actionName string, del func(*MediatePayload[T, K])) Mediator[T, K] {
	mtr := Mediator[T, K]{action: del, observable: obs, actionName: actionName}
	go mtr.Listener()

	return mtr

}

func (mtr *Mediator[T, K]) Mediate(msg T) (res K) {
	// go func(resp chan eventMessage[T, K]) {
	ctx := context.Background()
	ctx, close := context.WithTimeout(ctx, time.Second*3)
	defer close()

	request := mtr.observable.EmitWithResponse(mtr.actionName, msg)
	resp := mtr.observable.Subscriber(request.CorrelationId.String())

	// <-resp
	defer mtr.observable.RemoveSitter(request.CorrelationId.String(), resp)
	select {
	case result := (<-resp):
		res = result.response

	case <-ctx.Done():
		fmt.Println("ctx timeout")
	}
	// }(resp)
	return res

}
func (mtr *Mediator[T, K]) Listener() {
	var test K
	for {

		request := (<-mtr.observable.Subscriber(mtr.actionName))
		p := MediatePayload[T, K]{Payload: request.Args}

		mtr.action(&p)
		res := p.Response

		if res != returnZero(test) {
			mtr.observable.Response(request.CorrelationId.String(), res)

		}
		// if request.withresponse {
		// }

	}
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
