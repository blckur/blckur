package handlers

import (
    "github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
    engine.GET("/user", user_get)
}
