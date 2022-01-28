package mediator

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type Mediator[T Input, K Output] struct {
	actionChan       chan func(*MediatePayload[T, K])
	observable       *Observable[T, K]
	actionName       string
	cancelationToken chan bool
}
type MediatePayload[T Input, K Output] struct {
	Payload  T
	Response K
}
type IMediator[T Input, K Output] interface {
	Mediate(msg T) (res K)
	Close()
	Mediator(*MediatePayload[T, K])
}

func (obs *Observable[T, K]) NewMediator(actionName string) Mediator[T, K] {

	mtr := Mediator[T, K]{actionChan: make(chan func(*MediatePayload[T, K]), 1), observable: obs, actionName: actionName, cancelationToken: make(chan bool)}
	mtr.start()

	return mtr

}
func (mtr *Mediator[T, K]) AddOrUpdateCallback(del func(*MediatePayload[T, K])) {
	mtr.observable.addCallback(mtr.actionName, del)
	//mtr.observable.action = del
	//mtr.actionChan <- del
}

//del func(*MediatePayload[T, K]
func (mtr *Mediator[T, K]) start() {
	go mtr.listener()
	//need to create channel to listen to subscribe init
	//time.Sleep(1 * time.Second)
}
func (mtr *Mediator[T, K]) Close() {
	go func() {
		mtr.cancelationToken <- true
	}()
}
func (mtr *Mediator[T, K]) Mediate(msg T) (res K) {
	// go func(resp chan eventMessage[T, K]) {

	ctx := context.Background()
	//super critic in order to prevent memory leak
	wrp := newEventWrapper[T, K](msg, true)
	resp := *mtr.observable.Subscribe(wrp.CorrelationId.String())
	request := mtr.observable.EmitWithResponse(mtr.actionName, wrp)
	ctx, close := context.WithTimeout(ctx, time.Second*3)
	defer mtr.observable.RemoveRSitter(request.CorrelationId.String(), mtr.actionName, &resp)
	defer close()
	r := make(chan K, 0)
	go func() {
		select {

		case result := <-resp:

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

func (mtr *Mediator[T, K]) listener() {
	//var zeroValue K
	req := mtr.observable.Subscribe(mtr.actionName)
	for {
		select {
		case <-mtr.cancelationToken:
			{
				fmt.Println("Canceled")
				return
			}
		// case actions := <-mtr.actionChan:
		// 	{
		// 		mtr.observable.action = actions
		// 	}

		case request := <-*req:
			{
				go func() {

					p := MediatePayload[T, K]{Payload: request.Args}
					del := mtr.observable.getCallBack(mtr.actionName)
					if del != nil {
						for _, act := range *del {
							act(&p)
						}
					}

					// if mtr.observable.action != nil {
					// 	mtr.observable.action(&p)
					// }

					res := p.Response
					sres := fmt.Sprint(res)
					if len(sres) > 0 {
						mtr.observable.EmitResponse(request.CorrelationId.String(), res)
						//	defer mtr.observable.removeSubscriber(request.CorrelationId.String(),req)
					}
				}()
			}
		default:
			{
				//defer mtr.observable.RemoveSitter(mtr.actionName, req)
				//defer close(*req)

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
func (mtr *Mediator[T, K]) actionHandler(p *MediatePayload[T, K]) {
	fmt.Println(p)

}
