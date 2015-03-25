package account

import (
	"reflect"
)

const (
	OAUTH1 = 1
	OAUTH2 = 2
)

var (
	authTypes = map[string]int{}
	authRegistry = map[string]reflect.Type{}
	clientRegistry = map[string]reflect.Type{}
	handlers = []func(){}
)

func register(name string, typ int, auth interface{}, client interface{},
		handler func()) {
	authTypes[name] = typ
	authRegistry[name] = reflect.TypeOf(auth)
	clientRegistry[name] = reflect.TypeOf(client)
	handlers = append(handlers, handler)
}
