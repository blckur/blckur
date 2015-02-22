package handlers

import (
    "github.com/gin-gonic/gin"
)

type Account struct {
    Email string `json:"email" binding:"required"`
}

func account_get(c *gin.Context) {
    acct := &Account{
        Email: "zach.huff.386@gmail.com",
    }
    c.JSON(200, acct)
}
