package main

import (
    "github.com/blckur/blckur/handlers"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    handlers.Register(router)
    router.Run(":3000")
}
