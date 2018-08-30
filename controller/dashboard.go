package controller

import (
	"cryptobot/app"
	"cryptobot/services/rates"
)

type Dashboard struct {
	Rows        []Row              `json:"rows"`
	NewUserInfo *app.PublicAPIUser `json:"newUserInfo"`
}

type Row struct {
	rates.CurrencyRates
	Trade *app.PublicAPITrade `json:"trade"`
}
