package controller

import (
	"cryptobot/app"
	"cryptobot/conf"
	"cryptobot/helper"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *Controller) login(ctx *gin.Context) {
	authToken := ctx.Query("auth_token")
	if authToken == "" {
		ctx.Redirect(http.StatusFound, c.Config.FE.IJRP)
		ctx.Abort()
		return
	}
	var auth app.UserAuth
	err := c.App.DB.Where("auth_token = ? AND is_used = 0 AND expiring_at > ?", authToken, time.Now()).First(&auth).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			ctx.Redirect(http.StatusFound, c.Config.FE.IJRP)
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = c.App.DB.Table("user_auth").Where("auth_token = ?", authToken).Update("is_used", true).Error
	if err != nil {
		logrus.WithError(err).Error("can't update user_auth")
	}

	var u app.User
	if err := c.App.DB.Where("outer_id = ?", auth.UserID).First(&u).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "outer_id = ?" + err.Error(),
		})
		return
	}

	claims := MyCustomClaims{
		u.ID,
		u.Email,
		u.Language,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(c.Config.Auth.JwtTtl)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sToken, err := token.SignedString([]byte(c.Config.Auth.JWTSecret))
	// logrus.Info(sToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	userActionLog := &app.UserActionLog{
		Action:    "login",
		UserID:    u.ID,
		CreatedAt: time.Now(),
	}
	c.App.DB.Create(&userActionLog)

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  "session_token",
		Value: url.QueryEscape(sToken),
		Path:  "/",
	})
	ctx.Redirect(http.StatusFound, "/")
}

func (c *Controller) signUp(ctx *gin.Context) {
	var req map[string]struct {
		UserID   int64  `json:"userId"`
		Language string `json:"language"`
		Email    string `json:"email"`
		Country  string `json:"country"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}
	if _, ok := req["data"]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not founded \"data\"",
		})
		return
	}
	if req["data"].UserID == 0 || req["data"].Language == "" || req["data"].Email == "" || req["data"].Country == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "\"userId\" and \"language\" and \"email\" and \"country\" must be set",
		})
		return
	}
	demoBalance := c.Config.DemoBalance
	if demoBalance <= 0 {
		demoBalance = 1
	}
	u := &app.User{
		PublicAPIUser: app.PublicAPIUser{
			OuterID:    req["data"].UserID,
			Balance:    demoBalance,
			Email:      req["data"].Email,
			Mode:       "demo",
			Leverage:   c.Config.Leverages[0].Leverage,
			BotStarted: true,
			Country:    req["data"].Country,
		},
		Language:         req["data"].Language,
		DefaultMargin:    c.Config.Leverages[0].DailyProfit,
		TourVisited:      true,
		TodayBaseBalance: demoBalance,
		CreatedAt:        time.Now(),
	}
	u.CalculateLeverageWithoutRegress(c.Config.Leverages)
	if err := c.App.DB.Create(u).Error; err != nil {
		// If error duplicate, check if the data is already there
		// If yes do not return error
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry") {
			var count int
			c.App.DB.Table("users").
				Where("outer_id = ?", req["data"].UserID).
				Where("email = ?", req["data"].Email).
				Count(&count)
			if count == 0 {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	go c.addGTMEvent(int(req["data"].UserID), "lead", 0)
	ctx.JSON(http.StatusOK, gin.H{
		"response": "ok",
	})
}

func (c *Controller) signIn(ctx *gin.Context) {
	var req map[string]struct {
		UserID int64 `json:"userId"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body " + err.Error(),
		})
		return
	}
	if _, ok := req["data"]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found \"data\" field.",
		})
		return
	}
	if req["data"].UserID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "\"userId\" must be set",
		})
		return
	}

	var Count int
	if c.App.DB.Table("users").Where("outer_id = ?", req["data"].UserID).Count(&Count); Count == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("not found userId %d in database", req["data"].UserID),
		})
		return
	}
	// TODO: AuthToken must be unique.
	auth := &app.UserAuth{
		UserID:     req["data"].UserID,
		AuthToken:  helper.PseudoSha2(),
		CreatedAt:  time.Now(),
		ExpiringAt: time.Now().Add(time.Second * time.Duration(c.Config.Auth.AuthTokenTTL)),
		IsUsed:     false,
	}
	if err := c.App.DB.Create(auth).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"response": auth.AuthToken,
	})
}

func (c *Controller) deposit(ctx *gin.Context) {
	var req map[string]struct {
		UserID  int64   `json:"userId"`
		Deposit float64 `json:"deposit"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if _, ok := req["data"]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not founded \"data\" field.",
		})
		return
	}
	if req["data"].UserID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "\"userId\" must be set",
		})
		return
	}
	if req["data"].Deposit <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "deposit should be > 0",
		})
		return
	}
	if helper.GetPrecision(req["data"].Deposit) > 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "deposit precision should be >= 1e-8",
		})
		return
	}
	var u app.User
	if err := c.App.DB.Where("outer_id = ?", req["data"].UserID).First(&u).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	deposit := app.DepositLog{
		Amount:      req["data"].Deposit,
		UserID:      u.ID,
		IsFirstTime: false,
		CreatedAt:   time.Now(),
	}

	if u.Mode == "demo" {
		if err := c.Robots.StopBot(u.ID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": err.Error(),
			})
			return
		}

		u.Mode = "real"
		u.Balance = req["data"].Deposit
		u.Profit = 0
		u.Margin = 0
		u.TodayBaseBalance = req["data"].Deposit
		u.CalculateLeverageWithRegress(c.Config.Leverages)

		deposit.IsFirstTime = true

		if err := c.App.DB.Where("user = ? AND mode = ?", u.ID, "demo").Delete(app.Trade{}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if u.Mode == "real" {
		u.Balance += req["data"].Deposit
		u.TodayBaseBalance += req["data"].Deposit
		u.CalculateLeverageWithoutRegress(c.Config.Leverages)
	} else {
		logrus.Errorf("unexpected user mode %s. Aborting\n", u.Mode)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if err := c.App.DB.Save(&u).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.App.DB.Save(&deposit).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Add notification
	var Message string
	var EventName string
	// check if user already deposited
	var Count int
	c.App.DB.Table("notifications").Where("user = ?", int(req["data"].UserID)).Count(&Count)
	if Count > 0 {
		//do something here
		EventName = "deposit"
		Message = fmt.Sprintf("Your deposit of %v BTC was successful.", req["data"].Deposit)
	} else {
		EventName = "firstDeposit"
		Message = fmt.Sprintf("Congratulations! Your first deposit of %v BTC was successful.", req["data"].Deposit)
	}
	c.addGTMEvent(int(req["data"].UserID), EventName, req["data"].Deposit)
	c.addNotification(Message, int(req["data"].UserID))
	ctx.JSON(http.StatusOK, gin.H{
		"response": "ok",
	})
}

func (c *Controller) withdraw(ctx *gin.Context) {
	var req map[string]struct {
		UserID   int64   `json:"userId"`
		Withdraw float64 `json:"withdraw"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if _, ok := req["data"]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not founded \"data\" field.",
		})
		return
	}
	if req["data"].UserID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "\"userId\" must be set",
		})
		return
	}
	if req["data"].Withdraw <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "withdraw should be > 0",
		})
		return
	}
	if helper.GetPrecision(req["data"].Withdraw) > 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "withdraw precision should be >= 1e-8",
		})
		return
	}
	var u app.User
	if err := c.App.DB.Where("outer_id = ?", req["data"].UserID).First(&u).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	deposit := app.DepositLog{
		Amount:      req["data"].Withdraw * -1,
		UserID:      u.ID,
		IsFirstTime: false,
		CreatedAt:   time.Now(),
	}

	if u.Mode == "real" || u.Mode == "demo" {
		if err := c.Robots.StopBot(u.ID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": err.Error(),
			})
			return
		}
		u.Balance -= req["data"].Withdraw
		u.TodayBaseBalance -= req["data"].Withdraw
		u.Profit = 0
		u.CalculateLeverageWithoutRegress(c.Config.Leverages)
	} else {
		logrus.Errorf("unexpected user mode %s. Aborting\n", u.Mode)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if err := c.App.DB.Save(&deposit).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.App.DB.Save(&u).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	var Message string
	Message = fmt.Sprintf("Your withdrawal of %v BTC was success", req["data"].Withdraw)
	c.addNotification(Message, int(req["data"].UserID))
	c.addGTMEvent(int(req["data"].UserID), "withdraw", req["data"].Withdraw)
	ctx.JSON(http.StatusOK, gin.H{
		"response": "ok",
	})
}

func (c *Controller) getUser(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found userID parameter",
		})
		return
	}

	var u app.User
	if err := c.App.DB.Where("outer_id = ?", userID).First(&u).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		u.CalculateLeverageWithoutRegress(c.Config.Leverages)
		ctx.JSON(http.StatusOK, gin.H{
			"response": u,
		})
	}
}

func (c *Controller) getUsers(ctx *gin.Context) {
	var users []app.User
	if err := c.App.DB.Table("users").Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i := range users {
		users[i].CalculateLeverageWithoutRegress(c.Config.Leverages)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"response": users,
	})
}

func (c *Controller) changeMargin(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found userID parameter",
		})
		return
	}
	var req map[string]struct {
		Margin *float64 `json:"customDailyMargin"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if _, ok := req["data"]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not founded \"data\" field.",
		})
		return
	}

	if req["data"].Margin != nil && *req["data"].Margin < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "margin should be >= 0",
		})
		return
	}

	var count int
	if err := c.App.DB.Table("users").Where("outer_id = ?", userID).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if count < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found user with id = " + userID,
		})
		return
	}

	if err := c.App.DB.Table("users").Where("outer_id = ?", userID).Update("custom_margin", req["data"].Margin).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"response": "ok",
		})
	}
}

func (c *Controller) changeTradeBounds(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found userID parameter",
		})
		return
	}
	var req struct {
		DailyTradeBoundsLower  float64 `json:"dailyTradeBoundsLower"`
		DailyTradeBoundsUpper  float64 `json:"dailyTradeBoundsUpper"`
		ResultTradeBoundsLower float64 `json:"resultTradeBoundsLower"`
		ResultTradeBoundsUpper float64 `json:"resultTradeBoundsUpper"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int
	if err := c.App.DB.Table("users").Where("outer_id = ?", userID).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if count < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found user with id = " + userID,
		})
		return
	}

	if err := c.App.DB.Table("users").Where("outer_id = ?", userID).Update(&app.User{
		HasCustomBounds:              true,
		DailyTradeProfitBoundsLower:  req.DailyTradeBoundsLower,
		DailyTradeProfitBoundsUpper:  req.DailyTradeBoundsUpper,
		ResultTradeProfitBoundsLower: req.ResultTradeBoundsLower,
		ResultTradeProfitBoundsUpper: req.ResultTradeBoundsUpper,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"response": "ok",
		})
	}
}

func (c *Controller) resetTradeBounds(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found userID parameter",
		})
		return
	}

	var count int
	if err := c.App.DB.Table("users").Where("outer_id = ?", userID).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if count < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found user delete with id = " + userID,
		})
		return
	}

	query := `UPDATE users SET has_custom_bounds = false,
				daily_trade_profit_bounds_lower = '0',
				daily_trade_profit_bounds_upper = '0',
				result_trade_profit_bounds_upper = '0',
				result_trade_profit_bounds_lower = '0'
			  WHERE outer_id = ?`
	if err := c.App.DB.Exec(query, userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"response": "ok",
		})
	}
}

func (c *Controller) changeFixFactor(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found userID parameter",
		})
		return
	}
	var req struct {
		MinFixFactor float64 `json:"minFixFactor"`
		MaxFixFactor float64 `json:"maxFixFactor"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int
	if err := c.App.DB.Table("users").Where("outer_id = ?", userID).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if count < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found user with id = " + userID,
		})
		return
	}

	if err := c.App.DB.Table("users").Where("outer_id = ?", userID).Update(&app.User{
		HasCustomFixFactor: true,
		MinFixFactor:       req.MinFixFactor,
		MaxFixFactor:       req.MaxFixFactor,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"response": "ok",
		})
	}
}

func (c *Controller) resetFixFactor(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found userID parameter",
		})
		return
	}

	var count int
	if err := c.App.DB.Table("users").Where("outer_id = ?", userID).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if count < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not found user delete with id = " + userID,
		})
		return
	}

	query := `UPDATE users SET has_custom_fix_factor = false,
				min_fix_factor = '0',
				max_fix_factor = '0',
			  WHERE outer_id = ?`
	if err := c.App.DB.Exec(query, userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"response": "ok",
		})
	}
}

func (c *Controller) getConfig(ctx *gin.Context) {
	type config struct {
		DemoBalance float64             `json:"demoBalance,omitempty"`
		Leverages   []conf.LeverageConf `json:"leverages"`
		MinBalance  float64             `json:"minBalance,omitempty"`

		DemoRateUpdateMinSeconds int `json:"demoRatesUpdateIntervalMinMS"`
		DemoRateUpdateMaxSeconds int `json:"demoRatesUpdateIntervalMaxMS"`

		RealRateUpdateMinSeconds int `json:"realRatesUpdateIntervalMinMS"`
		RealRateUpdateMaxSeconds int `json:"realRatesUpdateIntervalMaxMS"`

		FixFactor conf.FixFactor `json:"fixFactor"`
	}

	resp := make(map[string]config)

	resp["data"] = config{
		DemoBalance:              c.Config.DemoBalance,
		Leverages:                c.Config.Leverages,
		MinBalance:               c.Config.MinBalance,
		DemoRateUpdateMinSeconds: c.Config.Rates.DemoRateUpdateMinSeconds,
		DemoRateUpdateMaxSeconds: c.Config.Rates.DemoRateUpdateMaxSeconds,
		RealRateUpdateMinSeconds: c.Config.Rates.RealRateUpdateMinSeconds,
		RealRateUpdateMaxSeconds: c.Config.Rates.RealRateUpdateMaxSeconds,
		FixFactor:                c.Config.FixFactor,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"response": resp,
	})
}

func (c *Controller) setConfig(ctx *gin.Context) {
	var req map[string]struct {
		DemoBalance float64             `json:"demoBalance,omitempty"`
		Leverages   []conf.LeverageConf `json:"leverages,omitempty"`
		MinBalance  float64             `json:"minBalance,omitempty"`

		DemoRateUpdateMinSeconds int `json:"demoRatesUpdateIntervalMinMS,omitempty"`
		DemoRateUpdateMaxSeconds int `json:"demoRatesUpdateIntervalMaxMS,omitempty"`

		RealRateUpdateMinSeconds int `json:"realRatesUpdateIntervalMinMS,omitempty"`
		RealRateUpdateMaxSeconds int `json:"realRatesUpdateIntervalMaxMS,omitempty"`

		FixFactor conf.FixFactor `json:"fixFactor"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	logrus.Printf("Data: %v", req["data"])

	if req["data"].DemoBalance != 0 {
		c.Config.DemoBalance = req["data"].DemoBalance
	}

	if req["data"].Leverages != nil && len(req["data"].Leverages) > 0 {
		c.Config.Leverages = req["data"].Leverages
		c.Core.Leverages = req["data"].Leverages
	}

	if req["data"].MinBalance != 0 {
		c.Config.MinBalance = req["data"].MinBalance
	}

	c.Config.Rates.DemoRateUpdateMinSeconds = req["data"].DemoRateUpdateMinSeconds
	c.Config.Rates.DemoRateUpdateMaxSeconds = req["data"].DemoRateUpdateMaxSeconds
	c.Config.Rates.RealRateUpdateMinSeconds = req["data"].RealRateUpdateMinSeconds
	c.Config.Rates.RealRateUpdateMaxSeconds = req["data"].RealRateUpdateMaxSeconds

	// Update Rates service config
	logrus.Printf("[TB] Old rate config", c.Rates.Config)
	c.Rates.Config.DemoRateUpdateMinSeconds = req["data"].DemoRateUpdateMinSeconds
	c.Rates.Config.DemoRateUpdateMaxSeconds = req["data"].DemoRateUpdateMaxSeconds
	c.Rates.Config.RealRateUpdateMinSeconds = req["data"].RealRateUpdateMinSeconds
	c.Rates.Config.RealRateUpdateMaxSeconds = req["data"].RealRateUpdateMaxSeconds
	logrus.Printf("[TB] New rate config", c.Rates.Config)

	if req["data"].FixFactor != (conf.FixFactor{}) {
		c.Config.FixFactor = req["data"].FixFactor
	}

	if err := conf.Save("./conf/config.json", c.Config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"response": "ok",
	})
}
