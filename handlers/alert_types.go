package handlers

import (
	"github.com/gin-gonic/gin"
)

var (
	alertTypes = map[string][]AlertType{
		"gmail": []AlertType{
			AlertType{
				Label: "All new messages",
				Type: "all",
			},
			AlertType{
				Label: "New messages matching sender",
				Type: "from",
				ValueType: "input",
				ValueLabel: "Enter complete or partial email address " +
					"of sender to match",
				ValueHolder: "Email address",
			},
			AlertType{
				Label: "New messages matching subject",
				Type: "subject",
				ValueType: "input",
				ValueLabel: "Enter search term to match in email subject",
				ValueHolder: "Search term",
			},
			AlertType{
				Label: "New messages matching message body",
				Type: "body",
				ValueType: "input",
				ValueLabel: "Enter search term to match in email body",
				ValueHolder: "Search term",
			},
		},
	}
)

type AlertType struct {
	Label string `json:"label"`
	Type string `json:"type"`
	ValueType string `json:"value_type"`
	ValueLabel string `json:"value_label"`
	ValueHolder string `json:"value_holder"`
}

func alertTypesGet(c *gin.Context) {
	c.JSON(200, types)
}
