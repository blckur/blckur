package account

import (
	"reflect"
	"sort"
)

const (
	OAUTH1 = 1
	OAUTH2 = 2
)

var (
	authTypes = map[string]int{}
	authRegistry = map[string]reflect.Type{}
	clientRegistry = map[string]reflect.Type{}
	accountTypes = []*AccountType{}
	AlertTypes = map[string][]*AlertType{}
	alertLabels = map[string]map[string]string{}
	handlers = []func(){}
)

func sortAccountAlerts() {
	srt := &accountTypesSort{
		AccountTypes: accountTypes,
	}
	sort.Sort(srt)
}

func register(name string, label string, typ int, auth interface{},
		client interface{}, types []*AlertType, handler func()) {
	authTypes[name] = typ
	authRegistry[name] = reflect.TypeOf(auth)
	clientRegistry[name] = reflect.TypeOf(client)

	accountTypes = append(accountTypes, &AccountType{
		Label: label,
		Type: name,
	})
	sortAccountAlerts()
	for i, acctType := range accountTypes {
		acctType.Id = i
	}

	labels := map[string]string{}
	for i, typ := range types {
		typ.Id = i
		labels[typ.Type] = typ.Label
	}
	AlertTypes[name] = types
	alertLabels[name] = labels

	handlers = append(handlers, handler)
}
