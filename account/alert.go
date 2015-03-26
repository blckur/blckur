package account

import (
	"labix.org/v2/mgo/bson"
)

type AlertType struct {
	Id int `json:"id"`
	Label string `json:"label"`
	Type string `json:"type"`
	ValueType string `json:"value_type"`
	ValueLabel string `json:"value_label"`
	ValueHolder string `json:"value_holder"`
}

type Alert struct {
	Id bson.ObjectId `bson:"_id" json:"id"`
	Type string `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
	Label string `bson:"-" json:"label"`
}

func (a *Account) ParseAlerts() {
	for _, alrt := range a.Alerts {
		if alrt.Id == "" {
			alrt.Id = bson.NewObjectId()
		}
		alrt.Label = alertLabels[a.Type][alrt.Type]
	}

	a.sortAlerts()
}
