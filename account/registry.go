package account

import (
	"reflect"
	"sort"
)

var (
	authTypes      = map[string]int{}
	authRegistry   = map[string]reflect.Type{}
	clientRegistry = map[string]reflect.Type{}
	AccountTypes   = []*AccountType{}
	FilterTypes    = map[string][]*FilterType{}
	filterLabels   = map[string]map[string]string{}
	handlers       = []func(){}
)

func sortAccountTypes() {
	srt := &accountTypesSort{
		AccountTypes: AccountTypes,
	}
	sort.Sort(srt)
}

func Register(name string, label string, typ int, auth interface{},
	client interface{}, types []*FilterType, handler func()) {

	authTypes[name] = typ
	authRegistry[name] = reflect.TypeOf(auth)
	clientRegistry[name] = reflect.TypeOf(client)

	AccountTypes = append(AccountTypes, &AccountType{
		Label: label,
		Type:  name,
	})
	sortAccountTypes()
	for i, acctType := range AccountTypes {
		acctType.Id = i
	}

	labels := map[string]string{}
	for i, typ := range types {
		typ.Id = i
		labels[typ.Type] = typ.Label
	}
	FilterTypes[name] = types
	filterLabels[name] = labels

	handlers = append(handlers, handler)
}
