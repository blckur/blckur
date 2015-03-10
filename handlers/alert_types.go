package handlers

import (
	"github.com/gin-gonic/gin"
)

var (
	alertTypes = map[string][]AlertType{
		"gmail": []AlertType{
			AlertType{
				Label: "All new messages",
				Type: "all_email",
				ValueType: "",
			},
			AlertType{
				Label: "New messages matching filter",
				Type: "filter_email",
				ValueType: "gmail_filter",
			},
		},
	}
)

type AlertType struct {
	Label string `json:"label"`
	Type string `json:"type"`
	Value string `json:"value_type"`
}

func alertTypesGet(c *gin.Context) {
	c.JSON(200, types)
}
