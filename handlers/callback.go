package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/gin-gonic/gin"
	"github.com/blckur/blckur/utils"
)

func callbackTwitterGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	params := utils.ParseParams(c.Request)
	token := params.GetByName("oauth_token")
	verifier := params.GetByName("oauth_verifier")

	if token == "" || verifier == "" {
		c.AbortWithStatus(400)
		return
	}

	acct, err := account.NewTwitter(db, token, verifier)
	if err != nil {
		panic(err)
	}

	c.JSON(200, acct)
}
