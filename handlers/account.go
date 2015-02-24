package handlers

import (
    "github.com/gin-gonic/gin"
    "labix.org/v2/mgo/bson"
)

type Account struct {
    Id bson.ObjectId `json:"id" binding:"required"`
    UserId bson.ObjectId `json:"user_id" binding:"required"`
    Type string `json:"type" binding:"required"`
    Identity string `json:"identity" binding:"required"`
}

type AccountType struct {
    Label string `json:"label" binding:"required"`
    Type string `json:"type" binding:"required"`
}

func accountsTypesGet(c *gin.Context) {
    types := []*AccountType{
        &AccountType{
            Label: "Twitter",
            Type: "twitter",
        },
        &AccountType{
            Label: "GitHub",
            Type: "github",
        },
        &AccountType{
            Label: "Gmail",
            Type: "gmail",
        },
    }

    c.JSON(200, types)
}

func accountsGet(c *gin.Context) {
    accounts := []*Account{
        &Account{
            Id: bson.NewObjectId(),
            UserId: bson.NewObjectId(),
            Type: "gmail",
            Identity: "zach.huff.386@gmail.com",
        },
        &Account{
            Id: bson.NewObjectId(),
            UserId: bson.NewObjectId(),
            Type: "twitter",
            Identity: "@zachhuff386",
        },
        &Account{
            Id: bson.NewObjectId(),
            UserId: bson.NewObjectId(),
            Type: "github",
            Identity: "zachhuff386",
        },
    }

    c.JSON(200, accounts)
}
