package handlers

import (
    "github.com/gin-gonic/gin"
    "labix.org/v2/mgo/bson"
)

type User struct {
    Id bson.ObjectId `json:"id" binding:"required"`
    Email string `json:"email" binding:"required"`
}

func user_get(c *gin.Context) {
    acct := &User{
        Id: bson.NewObjectId(),
        Email: "zach.huff.386@gmail.com",
    }
    c.JSON(200, acct)
}
