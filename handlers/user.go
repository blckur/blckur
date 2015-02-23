package handlers

import (
    "github.com/gin-gonic/gin"
    "labix.org/v2/mgo/bson"
)

type User struct {
    Id bson.ObjectId `json:"id" binding:"required"`
    Email string `json:"email" binding:"required"`
}

var user = &User{
    Id: bson.NewObjectId(),
    Email: "zach.huff.386@gmail.com",
}

func userGet(c *gin.Context) {
    c.JSON(200, user)
}
