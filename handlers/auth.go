package handlers

import (
    "github.com/blckur/blckur/database"
    "github.com/blckur/blckur/user"
    "github.com/blckur/blckur/auth"
    "github.com/gin-gonic/gin"
)

type AuthData struct {
    Email string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func loginPost(c *gin.Context) {
    db := c.MustGet("db").(*database.Database)

    cook, err := auth.GetCookie(c)
    if err != nil {
        panic(err)
    }

    sess, err := cook.GetSession(db)
    switch err.(type) {
    case nil:
        c.JSON(200, sess)
        return
    case *auth.NotFoundError:
    default:
        panic(err)
    }

    data := &AuthData{}
    c.Bind(data)

    usr, err := user.FindUser(db, data.Email)
    switch err.(type) {
    case nil:
    case *user.NotFoundError:
        c.JSON(401, &ErrorData{
            Error: "auth_email_invalid",
            Message: "Email is invalid",
        })
        return
    default:
        panic(err)
    }

    if auth := usr.CheckPassword(data.Password); auth != true {
        c.JSON(401, &ErrorData{
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

func signupPost(c *gin.Context) {
    db := c.MustGet("db").(*database.Database)
    data := &AuthData{}
    c.Bind(data)

    usr, err := user.NewUser(db, data.Email, data.Password)
    switch err.(type) {
        case nil:
        case *user.ExistsError:
            c.JSON(401, &ErrorData{
                Error: "signup_email_exists",
                Message: "Email is already signed up",
            })
            return
        default:
            panic(err)
    }

    cook, err := auth.GetCookie(c)
    if err != nil {
        panic(err)
    }

    _, err = cook.NewSession(db, usr.Id)
    if err != nil {
        panic(err)
    }

    c.JSON(200, usr)
}
