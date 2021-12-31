package mediator

import (
	"reflect"
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

func (mtr *Mediator[T, K]) Mediate(msg T) K {
	// go func(resp chan eventMessage[T, K]) {
	request := mtr.observable.EmitWithResponse(mtr.actionName, msg)
	resp := mtr.observable.Subscriber(request.CorrelationId.String())

	// <-resp
	defer mtr.observable.RemoveSitter(request.CorrelationId.String(), resp)

	// }(resp)
	return (<-resp).response

}
func (mtr *Mediator[T, K]) Listener() {
	for {

		request := (<-mtr.observable.Subscriber(mtr.actionName))
		p := MediatePayload[T, K]{Payload: request.Args}

		mtr.action(&p)
		res := p.Response

		if request.withresponse {
			mtr.observable.Response(request.CorrelationId.String(), res)
		}

	}
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
