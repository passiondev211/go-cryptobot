package robot

import (
	"cryptobot/app"
	"cryptobot/conf"
	"cryptobot/services/rates"
	"time"
)

type Robot struct {
	userID   int64
	finished chan bool
}

func Start(userID int64, ratesService *rates.Service, tradeCn chan app.Trade, config conf.Robot) *Robot {
	r := &Robot{
		userID:   userID,
		finished: make(chan bool, 1),
	}

	go func() {
		// defer func() { logrus.Debugf("robot %d killed", userID) }()
		// logrus.Debugf("robot %d started", userID)
		for {
			select {
			case <-r.finished:
				return
			default:
				for _, rate := range ratesService.GetRates("real") {
					if rate.BuyPrice < rate.SellPrice {
						time.Sleep(time.Millisecond * time.Duration(config.ExecTime))
						tradeCn <- app.Trade{
							Currency:     rate.Currency,
							BuyPrice:     rate.BuyPrice,
							SellPrice:    rate.SellPrice,
							BuyExchange:  rate.BuyExchange,
							SellExchange: rate.SellExchange,
							User:         userID,
						}
						// logrus.Debugf("robot %d sent create request", userID)
						time.Sleep(time.Second * time.Duration(config.TimeoutAfterNewTrade))
						break
					}
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	return r
}

func (r *Robot) Stop() {
	r.finished <- true
}
