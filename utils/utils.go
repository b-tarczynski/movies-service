package utils

import (
	"github.com/BarTar213/movies-service/models"
	"github.com/gin-gonic/gin"
)

//returns account information for user
//should be only used in handler functions that using middleware CheckAccount function
func GetAccount(c *gin.Context) *models.AccountInfo {
	account := c.Keys["account"].(models.AccountInfo)

	return &account
}
