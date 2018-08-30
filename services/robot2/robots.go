package robot2

import (
	"cryptobot/app"
	"cryptobot/conf"
	"cryptobot/services/rates"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	mx           *sync.RWMutex
	ratesService *rates.Service
	botsStates   map[int64]bool
	conf         conf.Robot
	application  *app.App
	tradeCn      chan app.BatchRobotTrade
}

func New(tradeCn chan app.BatchRobotTrade, config conf.Robot, ratesService *rates.Service, a *app.App) *Service {
	return &Service{
		mx:           &sync.RWMutex{},
		ratesService: ratesService,
		botsStates:   make(map[int64]bool),
		conf:         config,
		tradeCn:      tradeCn,
		application:  a,
	}
}

func (s *Service) Start() {
	//starting bots
	var users []app.User
	if err := s.application.DB.Table("users").Where("bot_started = true").Find(&users).Error; err != nil {
		logrus.Warningf("Cannot fetch users with active bots: %v", err)
	} else {
		logrus.Printf("Starting %v bots\n", len(users))
		s.mx.Lock()
		for _, user := range users {
			s.botsStates[user.ID] = true
		}
		s.mx.Unlock()
	}

	for w := 1; w <= 2; w++ {
		go s.worker(w)
	}
}

func (s *Service) StartBot(id int64) error {
	s.mx.Lock()
	s.botsStates[id] = true
	s.mx.Unlock()
	if err := s.application.DB.Table("users").Where("id = ?", id).Update("bot_started", true).Error; err != nil {
		logrus.Errorf("cannot update bot %v state in DB", id)
		return err
	}
	return nil
}

func (s *Service) StopBot(id int64) error {
	s.mx.Lock()
	s.botsStates[id] = false
	s.mx.Unlock()

	if err := s.application.DB.Table("users").Where("id = ?", id).Update("bot_started", false).Error; err != nil {
		logrus.Errorf("cannot update bot %v state in DB", id)
		return err
	}
	return nil
}

func (s *Service) BotState(id int64) (bool, error) {
	s.mx.RLock()
	b, ok := s.botsStates[id]
	s.mx.RUnlock()
	if ok {
		return b, nil
	}
	return false, fmt.Errorf("No such bot")
}

func (s *Service) worker(w int) {
	for mode := range s.ratesService.RatesChan {
		logrus.Printf("[TB] Receive trade message in worker %d", w)
		time.Sleep(time.Millisecond * time.Duration(s.conf.ExecTime))
		activeBots := []int64{}
		s.mx.RLock()
		for userID, state := range s.botsStates {
			if state {
				activeBots = append(activeBots, userID)
			}
		}
		s.mx.RUnlock()

		s.tradeCn <- app.BatchRobotTrade{
			TradeInfo: app.Trade{
				Mode: mode,
			},
			UserIDs: activeBots,
		}
		logrus.Debugf("%d robots  sent create request", len(activeBots))
	}
}
