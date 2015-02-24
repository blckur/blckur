package handlers

import (
    "github.com/gin-gonic/gin"
)

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
