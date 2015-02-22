package handlers

import (
    "github.com/gin-gonic/gin"
    "labix.org/v2/mgo/bson"
)

type Account struct {
    Id bson.ObjectId `json:"id" binding:"required"`
    Email string `json:"email" binding:"required"`
}

func account_get(c *gin.Context) {
    acct := &Account{
        Id: bson.NewObjectId(),
        Email: "zach.huff.386@gmail.com",
    }
    c.JSON(200, acct)
}
