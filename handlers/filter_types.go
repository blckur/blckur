package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/gin-gonic/gin"
)

func filterTypesGet(c *gin.Context) {
	acctType := c.Params.ByName("type")
	c.JSON(200, account.FilterTypes[acctType])
}
