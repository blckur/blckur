package handlers

import (
    "github.com/blckur/blckur/database"
    "github.com/blckur/blckur/auth"
    "github.com/gin-gonic/gin"
    "net/http"
    "io/ioutil"
)

func Limiter(c *gin.Context) {
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
}

func Database(c *gin.Context) {
    c.Set("db", database.GetDatabase())
}

func Session(required bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        db := c.MustGet("db").(*database.Database)

        cook, err := auth.GetCookie(c)
        if err != nil {
            panic(err)
        }

        sess, err := cook.GetSession(db)
        switch err.(type) {
            case nil:
            case *auth.NotFoundError:
                sess = nil
                err = nil
            default:
                panic(err)
        }

        if required && sess == nil {
            c.AbortWithStatus(401)
        }

        c.Set("session", sess)
    }
}

func Static(c *gin.Context) {
    resp, err := http.Get(
        "http://localhost:8080" + c.Params.ByName("path"))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    c.Data(200, resp.Header.Get("Content-Type"), body)
}

func Register(engine *gin.Engine) {
    engine.Use(Limiter)
    engine.Use(gin.Recovery())
    engine.Use(Database)

    engine.POST("/signup", signupPost)


    sessGroup := engine.Group("")
    sessGroup.Use(Session(false))

    sessGroup.POST("/login", loginPost)


    authGroup := engine.Group("")
    authGroup.Use(Session(true))

    authGroup.GET("/account_types", accountsTypesGet)

    authGroup.GET("/accounts", accountsGet)

    authGroup.GET("/user", userGet)
    authGroup.PUT("/user", userPut)

    authGroup.GET("/events", eventGet)

    engine.GET("/s/*path", Static)
}
