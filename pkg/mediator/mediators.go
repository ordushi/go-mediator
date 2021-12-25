package mediator

import (
	"fmt"
	"reflect"
)

type Mediator[T any] struct {
	action     func(T) interface{}
	observable *Observable[T]
	actionName string
}
type IMediator interface {
	Mediate()
	Forfit()
}

func (obs *Observable[T]) NewMediator(actionName string, del func(T) interface{}) Mediator[T] {
	mtr := Mediator[T]{action: del, observable: obs, actionName: actionName}
	return mtr

}
func (mtr *Mediator[T]) Mediate(msg T) interface{} {
	var resp chan eventMessage[interface{}]
	go func() {
		request := mtr.observable.EmitWithResponse(mtr.actionName, msg)
		resp = mtr.observable.Subscriber(request.CorrelationId.String())

	}()
	return (<-(resp)).Args

}
func (mtr *Mediator[T]) Listener() {

	for {

		request := (<-mtr.observable.Subscriber(mtr.actionName))
		res := mtr.action(request.Args.(T))
		if res != nil {
			request.response = res
			fmt.Println(res)
			mtr.observable.Emit(mtr.actionName, request.response)
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
