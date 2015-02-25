package handlers

import (
    "github.com/blckur/blckur/database"
    "github.com/blckur/blckur/user"
    "github.com/blckur/blckur/auth"
    "github.com/blckur/blckur/error"
    "github.com/gin-gonic/gin"
    "labix.org/v2/mgo"
)

type LoginData struct {
    Email string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func loginPost(c *gin.Context) {
    db := database.GetDatabase()

    cook, err := auth.GetCookie(c)
    if err != nil {
        panic(err)
    }

    sess, err := cook.GetSession(db)
    if err == nil {
        c.JSON(200, sess)
        return
    } else if err != mgo.ErrNotFound {
        panic(err)
    }

    data := &LoginData{}
    c.Bind(data)

    usr, err := user.FindUser(db, data.Email)
    if err != nil {
        if err == mgo.ErrNotFound {
            c.JSON(401, &error.ErrorData{
                Error: "auth_email_invalid",
                Message: "Email is invalid",
            })
            return
        }
        panic(err)
    }

    if auth := usr.CheckPassword(data.Password); auth != true {
        c.JSON(401, &error.ErrorData{
            Error: "auth_password_invalid",
            Message: "Password is invalid",
        })
        return
    }

    sess, err = cook.NewSession(db, usr.Id)
    if err != nil {
        panic(err)
    }

    c.JSON(200, sess)
}
