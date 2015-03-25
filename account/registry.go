package account

import (
	"reflect"
)

var (
	registry = map[string]reflect.Type{}
)

func register(name string, obj interface{}) {
	registry[name] = reflect.TypeOf(obj)
}
