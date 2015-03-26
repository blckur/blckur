package handlers

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/utils"
	"github.com/blckur/blckur/settings"
	"github.com/gin-gonic/gin"
)

func callbackGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	acctType := c.Params.ByName("type")

	auth, authType, err := account.GetAuth(acctType)
	if err != nil {
		c.JSON(400, &errorData{
			Error: "unknown_type",
			Message: "Unknown account type",
		})
	}

	params := utils.ParseParams(c.Request)
	var x string
	var y string
	var denied bool

	if authType == account.OAUTH1 {
		x = params.GetByName("oauth_token")
		y = params.GetByName("oauth_verifier")
		denied = params.GetByName("denied") != ""
	} else {
		x = params.GetByName("state")
		y = params.GetByName("code")

		error := params.GetByName("error")

		switch error {
		case "":
			denied = false
		case "access_denied":
			denied = true
		default:
			c.AbortWithStatus(400)
			return
		}
	}

	if x == "" || y == "" {
		c.AbortWithStatus(400)
		return
	}

	if !denied {
		_, err = auth.Authorize(db, x, y)
		if err != nil {
			panic(err)
		}
	}

	c.Redirect(301, settings.System.AppHome)
}
