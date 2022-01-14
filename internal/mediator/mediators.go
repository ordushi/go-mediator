package mediator

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type Mediator[T Input, K Output] struct {
	action           func(*MediatePayload[T, K])
	observable       *Observable[T, K]
	actionName       string
	cancelationToken chan bool
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

	mtr := Mediator[T, K]{action: del, observable: obs, actionName: actionName, cancelationToken: make(chan bool)}
	mtr.Start()

	return mtr

}
func (mtr *Mediator[T, K]) Start() {
	go mtr.Listener()
	//need to create channel to listen to subscribe init
	time.Sleep(1 * time.Second)
}
func (mtr *Mediator[T, K]) Stop() {
	go func() {
		mtr.cancelationToken <- true
	}()
}
func (mtr *Mediator[T, K]) Mediate(msg T) (res K) {
	// go func(resp chan eventMessage[T, K]) {

	ctx := context.Background()
	ctx, close := context.WithTimeout(ctx, time.Second*3)
	//super critic in order to prevent memory leak
	defer close()

	request := mtr.observable.EmitWithResponse(mtr.actionName, msg)
	resp := mtr.observable.Subscriber(request.CorrelationId.String())
	defer mtr.observable.RemoveRSitter(request.CorrelationId.String(), mtr.actionName, resp)
	r := make(chan K, 1)
	go func() {
		select {

		case result := <-*resp:
			r <- result.response

		case <-ctx.Done():
			var s K = request.response
			r <- s
			//timeout error maybe return error in the future
			fmt.Println("ctx timeout")
		}
	}()

	res = <-r
	return res

}
func (mtr *Mediator[T, K]) Listener() {
	var zeroValue K
	req := mtr.observable.Subscriber(mtr.actionName)
	for {
		select {
		case <-mtr.cancelationToken:
			{
				fmt.Println("Canceled")
				return
			}
		default:
			{
				//defer mtr.observable.RemoveSitter(mtr.actionName, req)
				//defer close(*req)
				request := <-*req
				p := MediatePayload[T, K]{Payload: request.Args}
				mtr.action(&p)
				res := p.Response
				if res != returnZero(zeroValue) {
					mtr.observable.Response(request.CorrelationId.String(), res)
				}
			}
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
