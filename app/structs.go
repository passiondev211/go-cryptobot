package app

import (
	"cryptobot/conf"
	"time"
)

type Leverage struct {
	MinBalance float64 `json:"minBalance"`
	Value      int     `json:"value"`
}

type PublicAPIUser struct {
	OuterID              int64     `json:"id"`
	Balance              float64   `json:"balance"`
	Margin               float64   `json:"margin" sql:"-"`
	Profit               float64   `json:"profit" gorm:"column:current_profit"`
	Mode                 string    `json:"mode"`
	BotStarted           bool      `json:"botStarted"`
	CurrentLeverageValue int       `json:"currentLeverageValue" sql:"-"`
	NextLeveragePub      *Leverage `json:"nextLeverage" sql:"-"`
	Leverage             int       `json:"leverage"`
	Email                string    `json:"email"`
	Country              string    `json:"country"`
}

// User includes model of information returns from /v1/user-info
// TODO(tiabc): Take another look at if this really `User`.
type User struct {
	PublicAPIUser
	ID int64 `json:"-"`

	DefaultMargin float64            `json:"-" sql:"-"`
	CustomMargin  *float64           `json:"-"`
	Language      string             `json:"-"`
	NextLeverage  *conf.LeverageConf `json:"-" sql:"-"`
	TourVisited   bool               `json:"-"`

	DailyTradeProfitBoundsLower  float64 `json:"-"`
	DailyTradeProfitBoundsUpper  float64 `json:"-"`
	ResultTradeProfitBoundsLower float64 `json:"-"`
	ResultTradeProfitBoundsUpper float64 `json:"-"`
	HasCustomBounds              bool    `json:"-"`

	MinFixFactor       float64 `json:"-"`
	MaxFixFactor       float64 `json:"-"`
	HasCustomFixFactor bool    `json:"-"`

	TodayBaseBalance float64   `json:"-"`
	CreatedAt        time.Time `json:"-"`

	ChatID string `json:"-"`
}

func (u *User) DailyTradeProfitBounds() TradeProfitBounds {
	return TradeProfitBounds{
		Lower: u.DailyTradeProfitBoundsLower,
		Upper: u.DailyTradeProfitBoundsUpper,
	}
}

func (u *User) ResultTradeProfitBounds() TradeProfitBounds {
	return TradeProfitBounds{
		Lower: u.ResultTradeProfitBoundsLower,
		Upper: u.ResultTradeProfitBoundsUpper,
	}
}

// TradeProfitBounds represent a lower and upper bound.
type TradeProfitBounds struct {
	Lower, Upper float64
}

type PublicAPITrade struct {
	ExpectedProfit float64 `json:"expectedProfit" sql:"-"`
	ActualProfit   float64 `json:"actualProfit" sql:"-"`
}

// Trade includes model of information returns from /v1/trades
type Trade struct {
	PublicAPITrade
	ID           uint64    `json:"-"`
	Currency     string    `json:"currency"`
	BuyPrice     float64   `json:"buyPrice"`
	SellPrice    float64   `json:"sellPrice"`
	BuyExchange  string    `json:"exchangeName1"`
	SellExchange string    `json:"exchangeName2"`
	Amount       float64   `json:"amount"`
	CreatedAt    time.Time `json:"createdAt"`
	User         int64     `json:"-"`
	Profit       float64   `json:"profit"`
	Mode         string    `json:"-"`
	FixFactor    float64   `json:"-"`
}

// CalculateProfitFromPricesAndAmount computes the profit
// from buy/sell prices and amount.
func (t Trade) CalculateProfitFromPricesAndAmount() float64 {
	return (t.SellPrice - t.BuyPrice) * t.Amount
}

type BatchRobotTrade struct {
	TradeInfo Trade
	UserIDs   []int64
}

// Stat includes model of information returns from /v1/trade-stats
type Stat struct {
	Day            time.Time `json:"date"`
	TradesCnt      int       `json:"tradesCount"`
	ProfitableRate float64   `json:"profitableRate"`
}

type UserAuth struct {
	UserID     int64
	AuthToken  string
	CreatedAt  time.Time
	ExpiringAt time.Time
	IsUsed     bool
}

func (UserAuth) TableName() string {
	return "user_auth"
}

type DepositLog struct {
	ID          int64
	Amount      float64
	UserID      int64
	IsFirstTime bool
	CreatedAt   time.Time
}

type HttpLog struct {
	ID           int64
	ClientIP     string
	URL          string
	Request      string
	Response     string
	CreatedAt    time.Time
	ResponseTime string
}

type Notifications struct {
	Id      float64
	Status  float64
	User    float64
	Message string
}

type GTMEvent struct {
	Id     float64
	User   float64
	Amount float64
	Event  string
}

type TradeLog struct {
	Id                   float64
	Mode                 string
	TotalUsers           int64
	SuccessfulUsers      int64
	ZeroFixFactorUsers   int64
	BelowMinBalanceUsers int64
	OutsideBoundUsers    int64
	StartAt              time.Time
	EndAt                time.Time
	DBTime               string
	TotalTime            string
}

type UserActionLog struct {
	Id        float64
	Action    string
	UserID    int64
	CreatedAt time.Time
}
