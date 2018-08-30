package controller

import (
	"cryptobot/app"
	"cryptobot/helper"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *Controller) exchangeRates(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"response": c.Rates.GetRates("real"),
	})
}

func (c *Controller) userInfo(ctx *gin.Context) {
	u, ok := userFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "not found user in context",
		})
		return
	}
	var err error
	u.CalculateLeverageWithoutRegress(c.Config.Leverages)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if u.Balance == 0 {
		u.Margin = 0
	} else {
		u.Margin = 100 * u.Profit / (u.Balance)
	}
	u.Balance = helper.Round(u.Balance, 8)
	ctx.JSON(http.StatusOK, gin.H{
		"response": u,
	})
}

func (c *Controller) tradeStats(ctx *gin.Context) {
	u, ok := userFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "not found user id in context",
		})
		return
	}

	var response []app.Stat

	if err := c.App.DB.Raw(`
		SELECT  COUNT(*) trades_cnt,
		SUM(buy_price >= trades.sell_price) / COUNT(*) nonprofitable_rate,  
		SUM(buy_price < trades.sell_price) / COUNT(*) profitable_rate,  
		DATE(created_at) day  
		FROM trades  
		WHERE user = ? AND mode = ? AND created_at <= ?
		GROUP BY DATE(created_at);`,
		u.ID, u.Mode, time.Now()).Find(&response).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

func (c *Controller) logout(ctx *gin.Context) {
	// u, _ := userFromContext(ctx)
	// err = c.App.DB.Table("user_auth").Where("auth_token = ?", authToken).Update("is_used", true).Error
	// if err != nil {
	// 	logrus.WithError(err).Error("can't update user_auth")
	// }
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  "session_token",
		Value: "",
		Path:  "/",
	})
	ctx.JSON(http.StatusOK, gin.H{
		"response": "goodby",
	})
}

func (c *Controller) retrieveTrades(ctx *gin.Context) {
	u, ok := userFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "no user in ctx",
		})
	}
	limitStr := ctx.Query("limit")

	limit, _ := strconv.Atoi(limitStr)
	if limit > 1000 {
		limit = 1000
	}
	if limit < 0 {
		limit = 0
	}

	var trs []app.Trade
	var afterLimit time.Time
	if u.Mode == "demo" {
		afterLimit = time.Now().Add(-3 * time.Hour)
	} else {
		afterLimit = time.Now().Add(-120 * time.Hour)
	}
	if err := c.App.DB.Limit(limit).Where(
		"user = ? AND mode = ? AND created_at <= ? AND created_at >= ?",
		u.ID,
		u.Mode,
		time.Now(),
		afterLimit,
	).Order("created_at desc").Find(&trs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"response": trs,
	})
}

func (c *Controller) createTrade(ctx *gin.Context) {
	newTrade, err := app.GetTradeFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// FIXME(tiabc): Why is the existence ignored?
	// FIXME(tiabc): Never use string keys in getting values from any context. Create constants for that.
	u, ok := userFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "not found user",
		})
		return
	}

	// TODO(tiabc): Write a short overview what it is.
	for _, cur := range c.Rates.GetRates(u.Mode) {
		if cur.Currency == newTrade.Currency {
			newTrade.BuyPrice = cur.BuyPrice
			newTrade.SellPrice = cur.SellPrice
			break
		}
	}
	newTrade.User = u.ID
	if u.Balance <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "insufficient funds",
		})
		return
	}

	if err := c.App.CreateTransaction(newTrade, u, c.Config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		// FIXME(tiabc): Try to avoid else branches. Just do `return` where needed (like above)
		// FIXME(tiabc): and leave the last instruction as the successful.
	} else {
		ctx.JSON(http.StatusCreated, gin.H{
			"response": newTrade,
		})
	}
}

func (c *Controller) appConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"response": c.Config.FE,
	})
}

func (c *Controller) startTrading(ctx *gin.Context) {
	u, ok := userFromContext(ctx)
	if !ok {
		return
	}
	var targetState bool
	limitStr := ctx.Query("state")

	if limitStr != "on" && limitStr != "off" {
		state, err := c.Robots.BotState(u.ID)
		if err != nil {
			targetState = true
		} else {
			targetState = !state
		}
	} else {
		if limitStr == "on" {
			targetState = true
		} else {
			targetState = false
		}
	}

	if !targetState {
		if err := c.Robots.StopBot(u.ID); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"response": "robot stopped",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	} else {
		if err := c.Robots.StartBot(u.ID); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"response": "robot started",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}
}

func (c *Controller) getTourVisited(ctx *gin.Context) {
	u, _ := userFromContext(ctx)

	logrus.Printf("Tour user: %v", u)

	ctx.JSON(http.StatusOK, gin.H{
		"response": u.TourVisited,
	})
}

func (c *Controller) setTourVisited(ctx *gin.Context) {
	u, _ := userFromContext(ctx)

	u.TourVisited = true
	if err := c.App.DB.Save(u).Error; err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"response": u.TourVisited,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func (c *Controller) dashboardInfo(ctx *gin.Context) {
	u, ok := userFromContext(ctx)
	if !ok {
		return
	}
	if dashboard, err := c.GatherUserInfo(u); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"response": dashboard,
		})
	}
}

func (c *Controller) GatherUserInfo(u *app.User) (*Dashboard, error) {
	var d Dashboard
	for _, rate := range c.Rates.GetRates(u.Mode) {
		r := Row{
			CurrencyRates: rate,
		}
		if lastTrade, ok := c.Core.GetLastUserTrade(u.ID); ok && lastTrade.Currency == rate.Currency {
			r.Trade = &lastTrade.PublicAPITrade
			r.Trade.ActualProfit = lastTrade.Profit
			c.Core.ClearLastTrade(u.ID) //WHY
		}
		d.Rows = append(d.Rows, r)
	}

	if state, err := c.Robots.BotState(u.ID); state && err != nil {
		u.BotStarted = true
	}

	if u.Balance == 0 {
		u.Margin = 0
	} else {
		u.Margin = 100 * u.Profit / (u.Balance)
	}
	iLastLeverage := len(c.Config.Leverages) - 1
	for i, lev := range c.Config.Leverages {
		if u.Balance >= lev.RequiredBalance {
			u.CurrentLeverageValue = lev.Leverage
			if i != iLastLeverage {
				u.NextLeveragePub = &app.Leverage{
					Value:      c.Config.Leverages[i+1].Leverage,
					MinBalance: c.Config.Leverages[i+1].RequiredBalance,
				}
			} else {
				u.NextLeverage = nil
			}
		}
	}

	d.NewUserInfo = &u.PublicAPIUser

	return &d, nil
}

// Notifications
func (c *Controller) addNotification(Message string, User int) error {
	err := c.App.DB.Exec("INSERT INTO notifications(message, status, user) VALUE(?, ?,?)", Message, 0, User).Error
	return err
}

func (c *Controller) getNotification(ctx *gin.Context) {
	u, _ := userFromContext(ctx)
	var result app.Notifications
	var Message string
	c.App.DB.Table("notifications").Where("status = ? AND user = ?", 0, u.OuterID).First(&result)
	if result.Id > 0 {
		if err := c.App.DB.Exec("UPDATE notifications SET status=? WHERE id = ?", 1, result.Id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		Message = result.Message
	}
	ctx.JSON(http.StatusOK, gin.H{
		"response": Message,
	})
}

// GTM Event
func (c *Controller) addGTMEvent(User int, Event string, Amount float64) error {
	err := c.App.DB.Exec("INSERT INTO gtm_events(user,event,amount) VALUE(?, ?, ?)", User, Event, Amount).Error
	return err
}

func (c *Controller) getGTMEvent(ctx *gin.Context) {
	u, _ := userFromContext(ctx)
	var result app.GTMEvent
	c.App.DB.Table("gtm_events").Where("user = ?", u.OuterID).First(&result)
	if result.Id > 0 {
		c.App.DB.Exec("DELETE FROM gtm_events WHERE id = ?", result.Id)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"response": result,
	})
}

func (c *Controller) getChatID(ctx *gin.Context) {
	u, _ := userFromContext(ctx)
	result := make(map[string]string)
	result["chat_id"] = u.ChatID
	ctx.JSON(http.StatusOK, gin.H{
		"response": result,
	})
}

func (c *Controller) setChatID(ctx *gin.Context) {
	u, _ := userFromContext(ctx)
	var req struct {
		ChatID string `json:"chat_id"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	u.ChatID = req.ChatID

	if err := c.App.DB.Save(&u).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"response": "ok",
	})
}
