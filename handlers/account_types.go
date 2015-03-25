package handlers

import (
	"github.com/gin-gonic/gin"
)

var (
	types = []accountType{
		accountType{
			Id: 1,
			Label: "Twitter",
			Type: "twitter",
		},
		accountType{
			Id: 2,
			Label: "GitHub",
			Type: "github",
		},
		accountType{
			Id: 3,
			Label: "Gmail",
			Type: "gmail",
		},
		accountType{
			Id: 4,
			Label: "HipChat",
			Type: "hipchat",
		},
		accountType{
			Id: 5,
			Label: "Bitly",
			Type: "bitly",
		},
		accountType{
			Id: 6,
			Label: "DigitalOcean",
			Type: "digitalocean",
		},
	}
)

type accountType struct {
	Id int `json:"id"`
	Label string `json:"label"`
	Type string `json:"type"`
}

func accountsTypesGet(c *gin.Context) {
	c.JSON(200, types)
}
