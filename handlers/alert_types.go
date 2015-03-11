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
				ValueType: "string",
				ValueLabel: "Enter complete or partial email address " +
					"of sender to match",
			},
			AlertType{
				Label: "New messages matching subject",
				Type: "subject",
				ValueType: "string",
				ValueLabel: "Enter search term to match in email subject",
			},
			AlertType{
				Label: "New messages matching message body",
				Type: "body",
				ValueType: "string",
				ValueLabel: "Enter search term to match in email body",
			},
		},
	}
)

type AlertType struct {
	Label string `json:"label"`
	Type string `json:"type"`
	ValueType string `json:"value_type"`
	ValueLabel string `json:"value_label"`
}

func alertTypesGet(c *gin.Context) {
	c.JSON(200, types)
}
