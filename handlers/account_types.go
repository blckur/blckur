package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/blckur/blckur/account"
)

func accountsTypesGet(c *gin.Context) {
	c.JSON(200, account.AccountTypes)
}
