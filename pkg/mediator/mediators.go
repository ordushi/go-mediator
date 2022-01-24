package mediator

import (
	"reflect"
)

var sngl map[reflect.Type]interface{} = make(map[reflect.Type]interface{})

func CreateOrGet[T any, K any]() Observable[T, K] {
	var t T
	x := reflect.TypeOf(t)
	var val interface{}
	ok := false
	if val, ok = sngl[x]; !ok {
		obs := newObservable[T, K]()
		sngl[x] = obs
		return obs

	}
	return val.(Observable[T, K])

}
func Remove[T any]() {
	var t T
	typ := reflect.TypeOf(t)
	delete(sngl, typ)
}
