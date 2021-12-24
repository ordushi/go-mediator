package mediator

import (
	"reflect"
	"time"
)

// type mediators[T Observable[T]] struct {
// 	observable T
// }

func new[T any]() Observable[T] {
	observ := Observable[T]{time: time.Now()}
	return observ
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

// func Get[T inter](x T) inter {
// 	return mediators[x]

// }
