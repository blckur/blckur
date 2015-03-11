package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/gin-gonic/gin"
)

var (
)
func alertTypesGet(c *gin.Context) {
	c.JSON(200, account.AlertTypes["gmail"])
}
