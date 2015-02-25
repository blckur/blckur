package main

import (
    "github.com/blckur/blckur/handlers"
    "github.com/blckur/blckur/database"
    "github.com/gin-gonic/gin"
)

func main() {
    err := database.Init()
    if err != nil {
        panic(err)
    }

    router := gin.Default()
    handlers.Register(router)
    router.Run(":3000")
}
