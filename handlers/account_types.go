package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/gin-gonic/gin"
)

func accountsTypesGet(c *gin.Context) {
	c.JSON(200, account.GetAccountTypes())
}
