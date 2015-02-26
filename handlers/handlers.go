package handlers

import (
    "github.com/blckur/blckur/database"
    "github.com/gin-gonic/gin"
    "net/http"
    "io/ioutil"
)

func Database(c *gin.Context) {
    c.Set("db", database.GetDatabase())
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
    engine.Use(gin.Recovery())
    engine.Use(Database)

    engine.POST("/login", loginPost)
    engine.POST("/signup", signupPost)

    engine.GET("/account_types", accountsTypesGet)

    engine.GET("/accounts", accountsGet)

    engine.GET("/user", userGet)
    engine.PUT("/user", userPut)

    engine.GET("/events", eventGet)

    engine.GET("/s/*path", Static)
}
