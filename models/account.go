package models

type AccountInfo struct {
	AccountId int    `header:"X-Account-Id" binding:"required"`
	Account   string `header:"X-Account" binding:"required"`
}
