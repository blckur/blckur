package handlers

import (
    "github.com/gin-gonic/gin"
)

func AccessControl() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
    }
}

func Register(engine *gin.Engine) {
    engine.Use(AccessControl())
    engine.GET("/user", user_get)
}
