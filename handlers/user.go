package handlers

import (
    "github.com/blckur/blckur/user"
    "github.com/blckur/blckur/error"
    "github.com/gin-gonic/gin"
    "labix.org/v2/mgo"
)

type LoginData struct {
    Email string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func userGet(c *gin.Context) {
}

func userPut(c *gin.Context) {
}

func loginPost(c *gin.Context) {
    data := &LoginData{}
    c.Bind(data)

    usr, err := user.FindUser(data.Email)
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

    c.JSON(200, usr)
}
