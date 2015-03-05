package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/utils"
	"github.com/blckur/blckur/settings"
	"github.com/gin-gonic/gin"
)

func callbackTwitterGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	params := utils.ParseParams(c.Request)
	token := params.GetByName("oauth_token")
	verifier := params.GetByName("oauth_verifier")
	denied := params.GetByName("denied")

	if denied == "" {
		if token == "" || verifier == "" {
			c.AbortWithStatus(400)
			return
		}

		_, err := account.NewTwitter(db, token, verifier)
		if err != nil {
			panic(err)
		}
	}

	c.Redirect(301, settings.System.Domain + "/s/") // TODO
}
