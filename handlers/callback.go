package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/gin-gonic/gin"
)

func callbackTwitterGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	params := c.Request.URL.Query()
	tokenList := params["oauth_token"]
	verifierList := params["oauth_verifier"]

	if len(tokenList) != 1 || len(verifierList) != 1 {
		c.AbortWithStatus(400)
		return
	}
	token := tokenList[0]
	verifier := verifierList[0]

	acct, err := account.NewTwitter(db, token, verifier)
	if err != nil {
		panic(err)
	}

	c.JSON(200, acct)
}
