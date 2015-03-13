package account

import (
	"labix.org/v2/mgo/bson"
)

var (
	AlertTypes = map[string][]AlertType{
		"gmail": []AlertType{
			AlertType{
				Id: 1,
				Label: "All new messages",
				Type: "all",
			},
			AlertType{
				Id: 2,
				Label: "New messages matching sender",
				Type: "from",
				ValueType: "input",
				ValueLabel: "Enter complete or partial email address " +
				"of sender to match",
				ValueHolder: "Email address",
			},
			AlertType{
				Id: 3,
				Label: "New messages matching subject",
				Type: "subject",
				ValueType: "input",
				ValueLabel: "Enter search term to match in email subject",
				ValueHolder: "Search term",
			},
			AlertType{
				Id: 4,
				Label: "New messages matching message body",
				Type: "body",
				ValueType: "input",
				ValueLabel: "Enter search term to match in email body",
				ValueHolder: "Search term",
			},
		},
	}
	AlertLabels = map[string]map[string]string{
		"gmail": map[string]string{
			"all": "All messages",
			"from": "Messages matching sender",
			"subject": "Messages matching subject",
			"body": "Messages matching message body",
		},
	}
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

func (a *Account) ParseEvents() {
	for _, alrt := range a.Alerts {
		if alrt.Id == "" {
			alrt.Id = bson.NewObjectId()
		}
		alrt.Label = AlertLabels[a.Type][alrt.Type]
	}
}
