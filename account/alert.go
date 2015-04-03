package account

import (
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo/bson"
)

type FilterType struct {
	Id          int    `json:"id"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	ValueType   string `json:"value_type"`
	ValueLabel  string `json:"value_label"`
	ValueHolder string `json:"value_holder"`
	ValueExample string `json:"value_example"`
}

type Filter struct {
	Id    bson.ObjectId `bson:"_id" json:"id"`
	Type  string        `bson:"type" json:"type"`
	Value string        `bson:"value" json:"value"`
	Label string        `bson:"-" json:"label"`
}

func (a *Account) ParseFilters() {
	filters := set.NewSet()

	for i, alrt := range a.Filters {
		key := alrt.Type + alrt.Value

		if filters.Contains(key) {
			a.Filters = append(a.Filters[:i], a.Filters[i+1:]...)
			continue
		}
		filters.Add(key)

		if alrt.Id == "" {
			alrt.Id = bson.NewObjectId()
		}
		alrt.Label = filterLabels[a.Type][alrt.Type]
	}

	a.sortFilters()
}
