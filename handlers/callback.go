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

		_, err := account.AuthTwitter(db, token, verifier)
		if err != nil {
			panic(err)
		}
	}

	c.Redirect(301, settings.System.AppHome)
}

func callbackGoogleGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	params := utils.ParseParams(c.Request)
	state := params.GetByName("state")
	code := params.GetByName("code")
	error := params.GetByName("error")

	if error == "" {
		if state == "" || code == "" {
			c.AbortWithStatus(400)
			return
		}

		_, err := account.NewGmail(db, state, code)
		if err != nil {
			panic(err)
		}
	} else if error != "access_denied" {
		c.AbortWithStatus(400)
		return
	}

	c.Redirect(301, settings.System.AppHome)
}
