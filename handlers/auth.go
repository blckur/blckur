package handlers

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/email"
	"github.com/blckur/blckur/session"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/user"
	"github.com/blckur/blckur/utils"
	"github.com/gin-gonic/gin"
	"net/mail"
)

type authData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
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
				Error:   "email_invalid",
				Message: "Email is invalid",
			})
			return
		default:
			c.AbortWithError(500, err)
			return
		}
	}

	if auth := usr.CheckPassword(data.Password); auth != true {
		c.JSON(401, &errorData{
			Error:   "password_invalid",
			Message: "Password is invalid",
		})
		return
	}

	cook := session.NewCookie(c)

	sess, err = cook.NewSession(db, usr.Id, data.Remember)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, sess)
}

func sessionDelete(c *gin.Context) {
	sess := c.MustGet("session").(*session.Session)

	err := sess.Remove()
	if err != nil {
		c.AbortWithError(500, err)
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
			Error:   "email_invalid",
			Message: "Email is invalid",
		})
		return
	}

	usr := user.User{
		Email: email,
	}

	err = usr.Initialize(db, data.Password)
	if err != nil {
		switch err.(type) {
		case *database.DuplicateKeyError:
			c.JSON(400, &errorData{
				Error:   "email_exists",
				Message: "Email is already signed up",
			})
		default:
			c.AbortWithError(500, err)
		}
		return
	}

	cook := session.NewCookie(c)

	_, err = cook.NewSession(db, usr.Id, true)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, usr)
}

func resetGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	key := c.Params.ByName("key")

	usr, err := user.GetSessionUser(db, key)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			c.AbortWithError(404, err)
		default:
			c.AbortWithError(500, err)
		}
		return
	}

	cook := session.NewCookie(c)

	_, err = cook.NewSession(db, usr.Id, true)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Redirect(302, "/#/reset")
}

func resetPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &authData{}
	c.Bind(data)

	usr, err := user.FindUser(db, data.Email)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			c.JSON(401, &errorData{
				Error:   "email_invalid",
				Message: "Email is invalid",
			})
			return
		default:
			c.AbortWithError(500, err)
			return
		}
	}

	key, err := usr.ResetPassword(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	client, err := email.New()
	if err != nil {
		return
	}
	defer client.Close()

	err = client.Send(&mail.Address{"", usr.Email}, "Blckur - Password reset",
		"You're receiving this email because you requested a password "+
			"reset for your user account at Blckur.\r\n\r\nPlease go to "+
			"the following page and choose a new password:\r\n\r\n"+
			settings.System.Domain+"/reset/"+key.Id+"\r\n\r\nThanks for "+
			"using our site!")
	if err != nil {
		return
	}

	c.JSON(200, usr)
}
