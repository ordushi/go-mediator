package mediator

import (
	"reflect"
)

type Mediator[T Input, K Output] struct {
	action     func(T) K
	observable *Observable[T, K]
	actionName string
}
type IMediator interface {
	Mediate()
	Forfit()
}

func (obs *Observable[T, K]) NewMediator(actionName string, del func(T) K) Mediator[T, K] {
	mtr := Mediator[T, K]{action: del, observable: obs, actionName: actionName}
	return mtr

}
func (mtr *Mediator[T, K]) Mediate(msg T) K {
	var resp chan eventMessage[T, K]
	go func() {
		request := mtr.observable.EmitWithResponse(mtr.actionName, msg)
		resp = mtr.observable.Subscriber(request.CorrelationId.String())

	}()
	return (<-(resp)).response

}
func (mtr *Mediator[T, K]) Listener() {
	for {

		request := (<-mtr.observable.Subscriber(mtr.actionName))
		res := mtr.action(request.Args)
		if request.withresponse {
			mtr.observable.Response(mtr.actionName, res)
		}

	}
}

func getType(myvar interface{}) string {
	if tpe := reflect.TypeOf(myvar); tpe.Kind() == reflect.Ptr {
		return "*" + tpe.Elem().Name()
	} else {
		return tpe.Name()
	}
}
