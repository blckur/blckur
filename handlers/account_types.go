package handlers

import (
	"github.com/gin-gonic/gin"
)

var (
	types = []*AccountType{
		&AccountType{
			Label: "Twitter",
			Type: "twitter",
		},
		&AccountType{
			Label: "GitHub",
			Type: "github",
		},
		&AccountType{
			Label: "Gmail",
			Type: "gmail",
		},
		&AccountType{
			Label: "HipChat",
			Type: "hipchat",
		},
		&AccountType{
			Label: "Bitly",
			Type: "bitly",
		},
		&AccountType{
			Label: "DigitalOcean",
			Type: "digitalocean",
		},
	}
)

type AccountType struct {
	Label string `json:"label" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func accountsTypesGet(c *gin.Context) {
	c.JSON(200, types)
}
