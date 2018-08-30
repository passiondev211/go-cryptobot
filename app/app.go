package app

import (
	"cryptobot/conf"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	// mysql driver for side effect.
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"strings"
	"time"

	"cryptobot/helper"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// App struct implemented domain layer.

type App struct {
	// Connect with percona database.
	DB *gorm.DB
}

func New(db *gorm.DB) *App {
	r := &App{
		DB: db,
	}

	return r
}

func (a *App) countTradesOfLastTimeout(userID int64, config conf.Main) int {
	interestTime := time.Now().Add(time.Second * -1)

	var c int
	err := a.DB.Table("trades").Where("user = ? AND created_at > ?", userID, interestTime).Count(&c).Error
	if err != nil {
		logrus.WithError(err).Error("failed get count of trades")
	}
	return c
}

func (a *App) CalcUserMargin(u *User) float64 {
	if u.Balance == 0 {
		return 0
	} else {
		return u.Profit / u.Balance
	}
}

func (a *App) CheckMargin(u *User, diffTradeProfit float64) error {
	return nil
}

func (a *App) GetProfitSumForPeriodFromNow(userId int64, period time.Duration) float64 {
	now := gorm.NowFunc()
	logrus.Debugf("app.App.GetProfitSumForPeriodFromNow: now: %v", now.String())
	since := now.Truncate(period)
	logrus.Debugf("app.App.GetProfitSumForPeriodFromNow: since: %v", since.String())

	type Result struct {
		Profit float64
	}
	var results []Result
	err := a.DB.Table("trades").
		Select("SUM(profit) as profit").
		Where("user = ? AND created_at >= ?", userId, since).
		Scan(&results).
		Error
	if err != nil {
		logrus.WithError(err).Error("app.App.GetProfitSumForPeriodFromNow failed")
		return 0
	}
	profit := results[0].Profit
	logrus.Debugf("app.App.GetProfitSumForPeriodFromNow: profit: %f", profit)
	return profit
}

func (a *App) GetTotalDepositForPeriodFromNow(userId int64, period time.Duration) float64 {
	now := gorm.NowFunc()
	logrus.Debugf("app.App.GetTotalDepositForPeriodFromNow: now: %v", now.String())
	since := now.Truncate(period)
	logrus.Debugf("app.App.GetTotalDepositForPeriodFromNow: since: %v", since.String())

	type Result struct {
		Deposit float64
	}
	var results []Result
	err := a.DB.Table("deposit_logs").
		Select("SUM(amount) as deposit").
		Where("user_id = ? AND created_at >= ? AND is_first_time = 0", userId, since).
		Scan(&results).
		Error
	if err != nil {
		logrus.WithError(err).Error("app.App.GetTotalDepositForPeriodFromNow failed")
		return 0
	}
	deposit := results[0].Deposit
	logrus.Debugf("app.App.GetTotalDepositForPeriodFromNow: deposit: %f", deposit)
	return deposit
}

func (a *App) GetCurrentPositionInPeriodFromNow(period time.Duration) float64 {
	now := gorm.NowFunc()
	since := now.Truncate(period)
	diff := now.Sub(since)
	result := diff.Seconds() / period.Seconds()
	return result
}

func (u *User) CalculateLeverageWithRegress(leverages []conf.LeverageConf) {
	u.Leverage = leverages[0].Leverage
	u.DefaultMargin = leverages[0].DailyProfit
	lastCount := len(leverages) - 1

	for i, lev := range leverages {
		if u.Balance >= lev.RequiredBalance {
			u.Leverage = lev.Leverage
			u.DefaultMargin = lev.DailyProfit
			if i != lastCount {
				u.NextLeverage = &leverages[i+1]
			} else {
				u.NextLeverage = nil
			}
		}
	}
}

func (u *User) CalculateLeverageWithoutRegress(leverages []conf.LeverageConf) {
	lastCount := len(leverages) - 1

	// Find current leverage
	for i := range leverages {
		if u.Balance >= leverages[i].RequiredBalance {
			u.DefaultMargin = leverages[i].DailyProfit
			if u.Leverage < leverages[i].Leverage {
				u.Leverage = leverages[i].Leverage
			}
			if i != lastCount {
				u.NextLeverage = &leverages[i+1]
			} else {
				u.NextLeverage = nil
			}
		} else {
			break
		}
	}
}

func (a *App) ResetUserProfit(u *User) error {
	u.Profit = 0
	return a.DB.Save(u).Error
}

func (a *App) CreateTransaction(t *Trade, u *User, config conf.Main) error {
	c := a.countTradesOfLastTimeout(t.User, config)
	// FIXME(tiabc): Let's get rid of Debug.TradesTimeout?
	if c > 0 {
		logrus.Error("maximum number of trades reached")
		return errors.New("maximum number of trades reached")
	}

	// TODO(tiabc): Write a comment why it's needed.
	t.Mode = u.Mode
	t.Amount = u.Balance * float64(u.Leverage)
	profit := (t.SellPrice - t.BuyPrice) * t.Amount
	u.CalculateLeverageWithoutRegress(config.Leverages)
	if err := a.CheckMargin(u, profit); err != nil {
		return err
	}

	t.Profit = helper.Round(profit, 8)
	u.Balance += t.Profit
	u.Profit += t.Profit
	u.CalculateLeverageWithoutRegress(config.Leverages)

	tx := a.DB.Begin()
	if a.DB.Error != nil {
		return a.DB.Error
	}
	if err := tx.Create(&t).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(&u).Error; err != nil {
		tx.Rollback()
		return err
	}
	// TODO(tiabc): What is it?
	if c != a.countTradesOfLastTimeout(t.User, config) {
		tx.Rollback()
		return errors.New("detected races, transaction canceled")
	}
	tx.Commit()
	return nil
}

func GetTradeFromContext(ctx *gin.Context) (*Trade, error) {
	var req map[string]Trade
	if err := ctx.BindJSON(&req); err != nil {
		return nil, err
	}
	if errs := req["data"].Validate(); len(errs) != 0 {
		return nil, errors.New(strings.Join(errs, ","))
	}
	tr, ok := req["data"]
	if !ok {
		return nil, errors.New("not founded \"data\" in reqquest body")
	}
	return &tr, nil
}

func (u *User) InitLeverages(levs []conf.LeverageConf) {
	if !u.HasCustomBounds {
		for _, lev := range levs {
			if u.Balance >= lev.RequiredBalance {
				u.DefaultMargin = lev.DailyProfit
				u.Leverage = lev.Leverage
				u.DailyTradeProfitBoundsLower = lev.DailyTradeBoundLower
				u.DailyTradeProfitBoundsUpper = lev.DailyTradeBoundUpper
				u.ResultTradeProfitBoundsLower = lev.ResultTradeBoundLower
				u.ResultTradeProfitBoundsUpper = lev.ResultTradeBoundUpper
			} else {
				break
			}
		}
	}
}
