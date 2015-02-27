package handlers

import (
	"github.com/blckur/blckur/auth"
	"github.com/blckur/blckur/utils"
	"github.com/gin-gonic/gin"
)

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func userGet(c *gin.Context) {
	sess := c.MustGet("session").(*auth.Session)

	usr, err := sess.GetUser()
	if err != nil {
		panic(err)
	}

	c.JSON(200, usr)
}

func userPut(c *gin.Context) {
	sess := c.MustGet("session").(*auth.Session)
	data := &UserData{}
	c.Bind(data)

	usr, err := sess.GetUser()
	if err != nil {
		panic(err)
	}

	if data.Email != "" {
		email, err := utils.ParseEmail(data.Email)
		if err != nil {
			c.JSON(400, &ErrorData{
				Error:   "email_invalid",
				Message: "Email is invalid",
			})
			return
		}

		usr.Email = email
	}

	if data.Password != "" {
		err = usr.SetPassword(data.Password)
		if err != nil {
			panic(err)
		}
	}

	err = usr.Commit()
	if err != nil {
		panic(err)
	}

	c.JSON(200, usr)
}
