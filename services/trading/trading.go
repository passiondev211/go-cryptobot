package trading

import (
	"cryptobot/app"
	"cryptobot/conf"
	"cryptobot/helper"
	"cryptobot/services/rates"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/remeh/sizedwaitgroup"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Service struct {
	app         *app.App
	NewTradesCh chan app.Trade
	Rates       *rates.Service
	Leverages   []conf.LeverageConf
	MinBalance  float64
	FixFactor   conf.FixFactor

	// LastTrade stores last user trade within 10 seconds
	LastTrade   map[int64]*app.Trade
	LastTradeMx *sync.RWMutex

	NewBatchRobotTradeCh chan app.BatchRobotTrade

	randGen *rand.Rand

	demoTradeInterval int
	realTradeInterval int

	alert *logrus.Logger
}

func New(a *app.App, rateService *rates.Service, leverages []conf.LeverageConf, minBalance float64, fixFactor conf.FixFactor, demoTradeInterval int, realTradeInterval int) *Service {
	alert := logrus.New()
	alert.Formatter = &logrus.JSONFormatter{}
	alert.Out = &lumberjack.Logger{
		Filename:   "log/alert.log",
		MaxSize:    50, // Megabytes
		MaxBackups: 7,
		MaxAge:     3, // days
	}

	s := &Service{
		NewTradesCh: make(chan app.Trade),
		app:         a,
		LastTrade:   make(map[int64]*app.Trade),
		LastTradeMx: &sync.RWMutex{},
		Rates:       rateService,
		Leverages:   leverages,
		MinBalance:  minBalance,
		FixFactor:   fixFactor,

		NewBatchRobotTradeCh: make(chan app.BatchRobotTrade),

		randGen: rand.New(rand.NewSource(time.Now().UnixNano())),

		demoTradeInterval: demoTradeInterval,
		realTradeInterval: realTradeInterval,

		alert: alert,
	}

	return s
}

func (s *Service) Start() {
	go func() {
		for {
			time.Sleep(time.Second * 1)
			toDelete := []int64{}
			s.LastTradeMx.RLock()
			for k, v := range s.LastTrade {
				if time.Since(v.CreatedAt) > time.Second*10 {
					toDelete = append(toDelete, k)
				}
			}
			s.LastTradeMx.RUnlock()

			s.LastTradeMx.Lock()
			for _, k := range toDelete {
				delete(s.LastTrade, k)
			}
			s.LastTradeMx.Unlock()
		}
	}()
	for trade := range s.NewTradesCh {
		u, err := s.initUser(trade.User)
		if err != nil {
			logrus.WithError(err).Errorf("can't init user (id=%d) during creating new trade", trade.User)
			continue
		}

		//logrus.Debugf("core recieved expected profit %.8f", (trade.SellPrice-trade.BuyPrice)*10)
		trade.Mode = u.Mode
		amountFactor := float64(rand.Int63n(100)+1) / 100
		trade.Amount = u.Balance * float64(u.Leverage) * amountFactor
		diff := trade.SellPrice - trade.BuyPrice
		temp := diff * trade.Amount
		//logrus.Debugf("trade amount %.8f diffPrices %.8f total %.8f", trade.Amount, diff, temp)
		trade.ExpectedProfit = temp
		//logrus.Debugf("core stored expected profit %.8f", trade.ExpectedProfit)
		//logrus.Debug("---")

		// replace expected prices to actual prices
		for _, rate := range s.Rates.GetRates(u.Mode) {
			if trade.Currency == rate.Currency {
				trade.BuyPrice = rate.BuyPrice
				trade.SellPrice = rate.SellPrice
				break
			}
		}

		profit := (trade.SellPrice - trade.BuyPrice) * trade.Amount
		if err := s.app.CheckMargin(u, profit); err != nil {
			// logrus.WithError(err).Errorf("invalid margin at user with outer id %d", trade.User)
			continue
		}

		trade.Profit = helper.Round(profit, 8)
		u.Balance += trade.Profit
		u.Profit += trade.Profit
		u.InitLeverages(s.Leverages)

		tx := s.app.DB.Begin()
		if s.app.DB.Error != nil {
			logrus.WithError(s.app.DB.Error).Error("failed save transaction")
			continue
		}
		if err := tx.Create(&trade).Error; err != nil {
			tx.Rollback()
			logrus.WithError(err).Error("failed save transaction")
			continue
		}
		if err := tx.Save(&u).Error; err != nil {
			tx.Rollback()
			logrus.WithError(err).Error("failed update user info")
			continue
		}
		tx.Commit()
		s.StoreTrade(trade)
	}
}

func (s *Service) StoreTrade(trade app.Trade) {
	s.LastTradeMx.Lock()
	s.LastTrade[trade.User] = &trade
	s.LastTradeMx.Unlock()
}

// ClearLastTrade Implemented save removing from trades cache
func (s *Service) ClearLastTrade(id int64) {
	s.LastTradeMx.Lock()
	if _, ok := s.LastTrade[id]; ok {
		delete(s.LastTrade, id)
	}
	s.LastTradeMx.Unlock()
}

func (s *Service) GetLastUserTrade(id int64) (app.Trade, bool) {
	s.LastTradeMx.RLock()
	defer s.LastTradeMx.RUnlock()
	r, ok := s.LastTrade[id]
	ret := app.Trade{}
	if ok {
		ret = *r

	}
	return ret, ok
}

func (s *Service) initUser(id int64) (*app.User, error) {
	var u app.User
	if err := s.app.DB.Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}
	if u.Balance <= 0 {
		return nil, errors.New("insufficient funds")
	}
	u.InitLeverages(s.Leverages)
	return &u, nil
}

/*
- Daily profit/loss Range:  -2% - 2%
- Result profit/loss Range:  1% - 2%

This means, during the day (Daily profit/loss Range)
the margin can be somewhere between plus and minus 2%.
However, at the end of the day the margin should be
positive and between 1% and 2%
(Result profit/loss Range).

Daily and result profit take into account the base balance
of today (balance at the end of last day + today's deposit - today's withdrawal).
To calculate the base balance for today we take
the current balance, and substract the profit for this
day so far. See `main.BaseBalanceJob`.

To smooth the rate we apply fix factor to both buy
and sell price. The fix factor is calculated based
on the balance of the previous day, current day
profit so far, new trade's profit and the boundaries.

If the daily margin is within bound, we apply no correction.math
If it's out of bounds - corrections need to be applied.

*/

func FloatLinearTransform(a, b float64, pos float64) float64 {
	return a + (b-a)*pos
}

func BoundsLinearTransform(a, b app.TradeProfitBounds, pos float64) app.TradeProfitBounds {
	return app.TradeProfitBounds{
		Lower: FloatLinearTransform(a.Lower, b.Lower, pos),
		Upper: FloatLinearTransform(a.Upper, b.Upper, pos),
	}
}

// var aFullDay = time.Second * 300
var aFullDay = time.Hour * 24

func initTradeFixInterval() {
	intv := os.Getenv("TRADE_FIX_FACTOR_INTERVAL")
	if intv != "" {
		val, err := time.ParseDuration(intv)
		if err != nil {
			panic(err)
		}
		aFullDay = val
	}
}

func init() {
	initTradeFixInterval()
}

type tradeFixMode uint8

const (
	tradeFixModeAbort   tradeFixMode = 1
	tradeFixModeCorrect tradeFixMode = 2
)

const currentTradeFixMode = tradeFixModeCorrect

func (s *Service) tradeFix(user *app.User, newTradeProfit float64, bounds app.TradeProfitBounds) float64 {
	profitForToday := user.Balance - user.TodayBaseBalance
	todayBaseBalance := user.TodayBaseBalance

	// ratioWithoutThisTrade := profitForToday / todayBaseBalance
	ratioWithThisTrade := (profitForToday + newTradeProfit) / todayBaseBalance

	if ratioWithThisTrade >= bounds.Lower && ratioWithThisTrade <= bounds.Upper {
		// We're within bound after the trade,
		// no need to adjust things.
		return 1
	}

	// We are out of bounds, we need to apply fix factor.

	// We first need to determine the goal.
	var goal float64
	if ratioWithThisTrade < bounds.Lower {
		goal = bounds.Lower
	} else {
		goal = bounds.Upper
	}

	// If the ratio with trade and the goal has different signs - we can't
	// possibly compensate that.
	// (Disable this for now)
	// if math.Signbit(goal) != math.Signbit(ratioWithThisTrade) {
	// 	// Just absorb the trade.
	// 	logrus.Debugf("trading.tradeFix: absorbing")
	// 	return 0
	// }

	switch currentTradeFixMode {
	case tradeFixModeAbort:
		// Instead of adjusting the trade values, we just abort the trade.
		var fixFactor float64 // = 0
		return fixFactor
	case tradeFixModeCorrect:
		// We calculate the factor from just fetching the values that we need.
		var fixFactor float64
		if newTradeProfit != 0 {
			fixFactor = ((goal * todayBaseBalance) - profitForToday) / newTradeProfit
		} else {
			fixFactor = 1
		}

		// Check ourselves here.
		// proof := (profitForToday + newTradeProfit*fixFactor) / todayBaseBalance

		return fixFactor
	default:
		// Unknown mode.
		panic("Unexpected trade fix mode")
	}
}

func (s *Service) getFixFactorBounds(user *app.User) app.TradeProfitBounds {
	pos := s.app.GetCurrentPositionInPeriodFromNow(aFullDay)
	//logrus.Debugf("trading.getFixFactorBounds: pos: %v", pos)
	bounds := BoundsLinearTransform(user.DailyTradeProfitBounds(), user.ResultTradeProfitBounds(), pos)
	//logrus.Debugf("trading.getFixFactorBounds: bounds: %v", bounds)
	return bounds
}

func (s *Service) BatchTradesHandler() {
	for batchTrade := range s.NewBatchRobotTradeCh {
		// logrus.Printf("batch: creating transactions for %v users\n", len(batchTrade.UserIDs))
		users := []app.User{}
		if err := s.app.DB.Where("bot_started = ? and mode = ? and balance >= ?", true, batchTrade.TradeInfo.Mode, s.MinBalance).Find(&users).Error; err != nil {
			logrus.WithError(err).Errorln("cannot fetch users for batch update")
		}
		tradeLog := &app.TradeLog{
			Mode:       batchTrade.TradeInfo.Mode,
			StartAt:    time.Now(),
			TotalUsers: int64(len(users)),
		}

		usersChunks := splitUsers(users, 1500)
		swg := sizedwaitgroup.New(8)
		for _, usersChunk := range usersChunks {
			swg.Add()
			go func(usersChunk []app.User) {
				defer swg.Done()
				preparedUsers := []app.User{}
				preparedTrades := []app.Trade{}

				for _, user := range usersChunk {
					currentRates := s.Rates.GetRates(user.Mode)
					i := rand.Intn(len(currentRates))
					newRates := currentRates[i]
					var randTime int
					if batchTrade.TradeInfo.Mode == "real" {
						randTime = rand.Intn(s.realTradeInterval)
					} else {
						randTime = rand.Intn(s.demoTradeInterval)
					}
					randDuration, err := time.ParseDuration(fmt.Sprintf("%ds", randTime))
					if err != nil {
						randDuration, _ = time.ParseDuration("1s")
					}

					trade := app.Trade{
						Currency:     newRates.Currency,
						BuyPrice:     newRates.BuyPrice,
						SellPrice:    newRates.SellPrice,
						BuyExchange:  newRates.BuyExchange,
						SellExchange: newRates.SellExchange,
						CreatedAt:    time.Now().Add(randDuration),
						Mode:         newRates.Mode,
					}

					user.InitLeverages(s.Leverages)
					trade.User = user.ID
					trade.Mode = user.Mode
					amountFactor := rand.Float64()
					if amountFactor == 0 {
						amountFactor = float64(rand.Int63n(100)+1) / 100.0
					}
					trade.Amount = user.Balance * float64(user.Leverage) * amountFactor
					diff := trade.SellPrice - trade.BuyPrice
					temp := diff * trade.Amount
					trade.ExpectedProfit = temp

					// Get fix facor.
					fixFactor := s.tradeFix(
						&user,
						trade.CalculateProfitFromPricesAndAmount(),
						s.getFixFactorBounds(&user),
					)

					// Appy fix factor.
					trade.FixFactor = fixFactor

					// Calculate profit with fix factor.
					profit := trade.CalculateProfitFromPricesAndAmount() * fixFactor

					if profit < 0 {
						//logrus.Debugf("BatchTradesHandler: lossy trade %v", trade)
					}

					minFixFactor := s.FixFactor.Min
					maxFixFactor := s.FixFactor.Max

					if user.HasCustomFixFactor {
						minFixFactor = user.MinFixFactor
						maxFixFactor = user.MaxFixFactor
					}

					// Only check fix factor range if previous profit != 0 or user is in real mode
					// In other word, demo user with 0 profit (usually new user with no trade yet)
					// will not check fix factor range
					if user.Profit != 0 || user.Mode == "real" {
						if fixFactor < minFixFactor || fixFactor > maxFixFactor {
							tradeLog.OutsideBoundUsers++
							continue
						}
					}

					if fixFactor == 0 {
						tradeLog.ZeroFixFactorUsers++
						continue
					}

					if err := s.app.CheckMargin(&user, profit); err != nil {
						continue
					}

					trade.Profit = helper.Round(profit, 8)

					finalBalance := user.Balance - trade.Profit
					if finalBalance < s.MinBalance {
						tradeLog.BelowMinBalanceUsers++
						continue
					}

					user.Balance += trade.Profit
					user.Profit += trade.Profit
					user.InitLeverages(s.Leverages)

					todaysProfitRatio := (user.Balance - user.TodayBaseBalance) / user.TodayBaseBalance
					if todaysProfitRatio < user.DailyTradeProfitBoundsLower || todaysProfitRatio > user.DailyTradeProfitBoundsUpper {
						s.alert.Errorf("BatchTradesHandler: Out of bound trade is done. User=%v, Email=%v, Margin=%v, Ratio:%v", user.ID, user.Email, user.Margin, todaysProfitRatio)
					}

					preparedTrades = append(preparedTrades, trade)
					preparedUsers = append(preparedUsers, user)
				}

				if len(preparedUsers) == 0 {
					return
				}
				tradeLog.SuccessfulUsers += int64(len(preparedUsers))

				tx := s.app.DB.Begin()
				if s.app.DB.Error != nil {
					logrus.WithError(s.app.DB.Error).Error("batch: failed creating transaction")
					tx.Rollback()
					return
				}
				f := true

				tradeQuery := `INSERT INTO trades (
					currency, buy_price, sell_price, buy_exchange, sell_exchange, amount,
					created_at, user, profit, mode, fix_factor
				) VALUES `
				valTempl := "('%s', %.10f, %.10f, '%s', '%s', %.10f, '%s', %d, %.10f, '%s', %.10f),"
				for _, trade := range preparedTrades {
					tradeQuery += fmt.Sprintf(
						valTempl, trade.Currency, trade.BuyPrice, trade.SellPrice,
						trade.BuyExchange, trade.SellExchange, trade.Amount,
						trade.CreatedAt.Format("2006-01-02 15:04:05"), trade.User,
						trade.Profit, trade.Mode, trade.FixFactor,
					)
					s.StoreTrade(trade)
				}
				tradeQuery = tradeQuery[:len(tradeQuery)-1]

				if err := tx.Exec(tradeQuery).Error; err != nil {
					logrus.WithError(err).Error("batch: failed save transaction")
					f = false
				}

				if !f {
					if tx != nil {
						tx.Rollback()
					}
					return
				}

				f = true
				for _, u := range preparedUsers {
					err := tx.Model(&u).Where("mode = ?", u.Mode).Update(app.User{
						PublicAPIUser: app.PublicAPIUser{
							Balance:  u.Balance,
							Profit:   u.Profit,
							Leverage: u.Leverage,
						},
						DailyTradeProfitBoundsLower:  u.DailyTradeProfitBoundsLower,
						DailyTradeProfitBoundsUpper:  u.DailyTradeProfitBoundsUpper,
						ResultTradeProfitBoundsLower: u.ResultTradeProfitBoundsLower,
						ResultTradeProfitBoundsUpper: u.ResultTradeProfitBoundsUpper,
					}).Error
					if err != nil {
						logrus.WithError(err).Error("batch: failed update user info")
						f = false
						break
					}
				}
				if f {
					tx.Commit()
				} else {
					if tx != nil {
						tx.Rollback()
					}
				}
			}(usersChunk)
		}
		swg.Wait()
		tradeLog.EndAt = time.Now()
		tradeLog.TotalTime = tradeLog.EndAt.Sub(tradeLog.StartAt).String()

		if err := s.app.DB.Create(&tradeLog).Error; err != nil {
			logrus.WithError(err).Error("batch: failed to save trade log")
		}
	}
}

func (s *Service) tradesBulkInsert(db *gorm.DB, unsavedRows []app.Trade) error {
	valueStrings := make([]string, 0, len(unsavedRows))
	valueArgs := make([]interface{}, 0, len(unsavedRows)*10)
	for _, post := range unsavedRows {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, post.Currency)
		valueArgs = append(valueArgs, post.BuyPrice)
		valueArgs = append(valueArgs, post.SellPrice)
		valueArgs = append(valueArgs, post.BuyExchange)
		valueArgs = append(valueArgs, post.SellExchange)
		valueArgs = append(valueArgs, post.Amount)
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, post.User)
		valueArgs = append(valueArgs, post.Profit)
		valueArgs = append(valueArgs, post.Mode)
	}
	stmt := fmt.Sprintf("INSERT INTO trades (currency, buy_price, sell_price, buy_exchange, sell_exchange, amount, created_at, user, profit, mode) VALUES %s",
		strings.Join(valueStrings, ","))
	err := db.Exec(stmt, valueArgs...).Error
	return err
}

func splitUsers(buffer []app.User, limit int) [][]app.User {
	var chunk []app.User
	chunks := make([][]app.User, 0, len(buffer)/limit+1)
	for len(buffer) >= limit {
		chunk, buffer = buffer[:limit], buffer[limit:]
		chunks = append(chunks, chunk)
	}

	if len(buffer) > 0 {
		chunks = append(chunks, buffer)
	}
	return chunks
}
