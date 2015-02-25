package handlers

import (
    "github.com/gin-gonic/gin"
)

func AccessControl() gin.HandlerFunc {
    return func(c *gin.Context) {
        headers := c.Writer.Header()

        headers.Add("Access-Control-Allow-Origin", "*")
        headers.Add("Access-Control-Allow-Methods",
            "GET,PUT,POST,DELETE")
        headers.Add("Access-Control-Max-Age", "43200")
        headers.Add("Access-Control-Allow-Headers",
            "Authorization,Content-Type,Accept,Origin,User-Agent,DNT," +
            "Cache-Control,X-Mx-ReqToken,Keep-Alive,X-Requested-With," +
            "If-Modified-Since")

        if (c.Request.Method == "OPTIONS") {
            c.AbortWithStatus(200)
        } else {
            c.Next()
        }
    }
}

func Register(engine *gin.Engine) {
    corsGroup := engine.Group("")
    corsGroup.Use(AccessControl())


    corsGroup.POST("/login", loginPost)
    corsGroup.OPTIONS("/login", nil)

    corsGroup.POST("/auth", authPost)
    corsGroup.OPTIONS("/auth", nil)

    corsGroup.GET("/account_types", accountsTypesGet)

    corsGroup.GET("/accounts", accountsGet)
    //corsGroup.GET("/accounts/detail", todo)

    corsGroup.GET("/user", userGet)
    corsGroup.PUT("/user", userPut)
    corsGroup.OPTIONS("/user", nil)

    engine.GET("/events", eventGet)
}
