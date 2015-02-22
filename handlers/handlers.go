package handlers

import (
    "github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
    engine.GET("/account", account_get)
}
