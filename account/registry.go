package account

import (
	"reflect"
)

var (
	authRegistry = map[string]reflect.Type{}
	clientRegistry = map[string]reflect.Type{}
	handlers = []func(){}
)

func register(name string, auth interface{}, client interface{},
		handler func()) {
	authRegistry[name] = reflect.TypeOf(auth)
	clientRegistry[name] = reflect.TypeOf(client)
	handlers = append(handlers, handler)
}
