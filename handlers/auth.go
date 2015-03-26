package handlers

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/session"
	"github.com/blckur/blckur/user"
	"github.com/blckur/blckur/utils"
	"github.com/gin-gonic/gin"
)

type authData struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Remember bool `json:"remember"`
}

func loginPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)
	data := &authData{}

	if !c.Bind(&data) {
		return
	}

	if sess != nil {
		c.JSON(200, sess)
		return
	}

	usr, err := user.FindUser(db, data.Email)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			c.JSON(401, &errorData{
				Error: "email_invalid",
				Message: "Email is invalid",
			})
			return
		default:
			panic(err)
		}
	}

	if auth := usr.CheckPassword(data.Password); auth != true {
		c.JSON(401, &errorData{
			Error: "password_invalid",
			Message: "Password is invalid",
		})
		return
	}

	cook, err := session.GetCookie(c)
	if err != nil {
		c.Fail(500, err)
		return
	}

	sess, err = cook.NewSession(db, usr.Id, data.Remember)
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, sess)
}

func sessionDelete(c *gin.Context) {
	sess := c.MustGet("session").(*session.Session)

	err := sess.Remove()
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.String(200, "")
}

func signupPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &authData{}
	c.Bind(data)

	email, err := utils.ParseEmail(data.Email)
	if err != nil {
		c.JSON(400, &errorData{
			Error: "email_invalid",
			Message: "Email is invalid",
		})
		return
	}

	usr := user.User{
		Email: email,
	}

	err = usr.Initialize(db, data.Password)
	switch err.(type) {
	case nil:
	case *database.DuplicateKeyError:
		c.JSON(400, &errorData{
			Error: "email_exists",
			Message: "Email is already signed up",
		})
		return
	default:
		c.Fail(500, err)
		return
	}

	cook, err := session.GetCookie(c)
	if err != nil {
		c.Fail(500, err)
		return
	}

	_, err = cook.NewSession(db, usr.Id, true)
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, usr)
}
