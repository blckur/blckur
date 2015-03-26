package account

type AccountType struct {
	Id int `json:"id"`
	Label string `json:"label"`
	Type string `json:"type"`
}

type accountTypesSort struct {
	AccountTypes []*AccountType
}

func (a *accountTypesSort) Len() int {
	return len(a.AccountTypes)
}

func (a *accountTypesSort) Less(i int, j int) bool {
	return a.AccountTypes[i].Label < a.AccountTypes[j].Label
}

func (a *accountTypesSort) Swap(i int, j int) {
	a.AccountTypes[i], a.AccountTypes[j] = a.AccountTypes[j], a.AccountTypes[i]
}

func GetAccountTypes() []*AccountType {
	return accountTypes
}
