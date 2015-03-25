package account

import (
	"reflect"
)

var (
	registry = map[string]reflect.Type{}
	handlers = []func(){}
)

func register(name string, obj interface{}, handler func()) {
	registry[name] = reflect.TypeOf(obj)
	handlers = append(handlers, handler)
}
