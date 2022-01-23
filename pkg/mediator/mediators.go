package mediator

import (
	"reflect"

	"ezkv.io/go-mediator/internal/mediator"
)

var sngl map[reflect.Type]interface{} = make(map[reflect.Type]interface{})

func AddOrGet[T any, K any]() mediator.Observable[T, K] {
	var t T
	x := reflect.TypeOf(t)
	var val interface{}
	ok := false
	if val, ok = sngl[x]; !ok {
		obs := mediator.New[T, K]()
		sngl[x] = obs
		return obs

	}
	return val.(mediator.Observable[T, K])

}
