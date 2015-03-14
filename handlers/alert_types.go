package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/gin-gonic/gin"
)

func alertTypesGet(c *gin.Context) {
	c.JSON(200, account.AlertTypes["gmail"])
}
