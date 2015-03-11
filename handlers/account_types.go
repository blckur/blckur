package handlers

import (
	"github.com/gin-gonic/gin"
)

var (
	types = []AccountType{
		AccountType{
			Id: 1,
			Label: "Twitter",
			Type: "twitter",
		},
		AccountType{
			Id: 2,
			Label: "GitHub",
			Type: "github",
		},
		AccountType{
			Id: 3,
			Label: "Gmail",
			Type: "gmail",
		},
		AccountType{
			Id: 4,
			Label: "HipChat",
			Type: "hipchat",
		},
		AccountType{
			Id: 5,
			Label: "Bitly",
			Type: "bitly",
		},
		AccountType{
			Id: 6,
			Label: "DigitalOcean",
			Type: "digitalocean",
		},
	}
)

type AccountType struct {
	Id int `json:"id"`
	Label string `json:"label"`
	Type string `json:"type"`
}

func accountsTypesGet(c *gin.Context) {
	c.JSON(200, types)
}
