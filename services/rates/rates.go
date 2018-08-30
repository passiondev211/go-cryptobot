package rates

import (
	"cryptobot/conf"
	"cryptobot/helper"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// All current real and modify rates
type Service struct {
	Config conf.Rates

	// Real rates price (update each 10 minutes).
	rates map[string]float64

	// Modified rates price (update each 2-5 seconds).
	modifiedDemoRates   []CurrencyRates
	modifiedDemoRatesMx *sync.RWMutex
	modifiedRealRates   []CurrencyRates
	modifiedRealRatesMx *sync.RWMutex
	RatesChan           chan string
	// DB implemented package link to database
	DB *gorm.DB

	// list of all available exchanges.
	exchanges []string

	// list of all available rates.
	currencies []string

	randGen *rand.Rand
}

// CurrencyRates include model of information returns from /v1/exchange-rates
type CurrencyRates struct {
	Currency     string  `json:"currency"`
	BuyPrice     float64 `json:"buyPrice"`
	SellPrice    float64 `json:"sellPrice"`
	BuyExchange  string  `json:"buyExchange"`
	SellExchange string  `json:"sellExchange"`
	Mode         string  `json:"-"`
}

func New(db *gorm.DB, config conf.Rates) *Service {
	s := Service{
		Config:              config,
		rates:               make(map[string]float64),
		modifiedDemoRatesMx: &sync.RWMutex{},
		modifiedRealRatesMx: &sync.RWMutex{},
		currencies:          initCurrencies(db),
		exchanges:           initExchanges(db),
		randGen:             rand.New(rand.NewSource(time.Now().UnixNano())),
		RatesChan:           make(chan string),
	}

	// Get prices from min-api.cryptocompare.com
	err := s.fetchRates()
	for err != nil {
		logrus.WithError(err).Error("Failed get initial rates. Trying again after ten seconds")
		time.Sleep(time.Second * 10)
		err = s.fetchRates()
	}

	// Make modified rates
	for cur, _ := range s.rates {
		s.modifiedDemoRates = append(s.modifiedDemoRates, s.makeRandomRate(cur, "demo"))
		s.modifiedRealRates = append(s.modifiedRealRates, s.makeRandomRate(cur, "real"))
	}

	return &s
}

// Run the service
func (s *Service) Start() {

	updateRatesTicker := time.After(s.actualRatesUpdateInterval())
	updateDemoModifiedRatesTicker := time.After(s.randomDemoRatesUpdateInterval())
	updateRealModifiedRatesTicker := time.After(s.randomRealRatesUpdateInterval())

	for {
		select {
		case <-updateRatesTicker:
			if err := s.fetchRates(); err != nil {
				logrus.WithError(err).Error("Failed update rates")
				updateRatesTicker = time.After(s.actualRatesUpdateInterval())
			}
		case <-updateDemoModifiedRatesTicker:
			logrus.Printf("[TB] Demo timer fired")
			s.modifiedDemoRatesMx.Lock()
			for i := 0; i < len(s.currencies); i++ {
				name := s.currencies[i]
				s.modifiedDemoRates[i] = s.makeRandomRate(name, "demo")
			}
			s.modifiedDemoRatesMx.Unlock()

			logrus.Printf("[TB] Sending demo trade message to robot service")
			select {
			case s.RatesChan <- "demo":
			default:
			}
			// fmt.Printf("new rates: %+v\n", s.modifiedRates[i])

			updateDemoModifiedRatesTicker = time.After(s.randomDemoRatesUpdateInterval())
		case <-updateRealModifiedRatesTicker:
			logrus.Printf("[TB] Real timer fired")
			s.modifiedRealRatesMx.Lock()
			for i := 0; i < len(s.currencies); i++ {
				name := s.currencies[i]
				s.modifiedDemoRates[i] = s.makeRandomRate(name, "real")
			}
			s.modifiedRealRatesMx.Unlock()

			logrus.Printf("[TB] Sending real trade message to robot service")
			select {
			case s.RatesChan <- "real":
			default:
			}
			// fmt.Printf("new rates: %+v\n", s.modifiedRates[i])

			updateRealModifiedRatesTicker = time.After(s.randomRealRatesUpdateInterval())
		}
	}
}

func (s *Service) actualRatesUpdateInterval() time.Duration {
	return time.Second * time.Duration(s.Config.RatesUpdate)
}

func (s *Service) randomDemoRatesUpdateInterval() time.Duration {
	t := s.randGen.Intn(1+s.Config.DemoRateUpdateMaxSeconds-s.Config.DemoRateUpdateMinSeconds) + s.Config.DemoRateUpdateMinSeconds
	return time.Second * time.Duration(t)
}

func (s *Service) randomRealRatesUpdateInterval() time.Duration {
	t := s.randGen.Intn(1+s.Config.RealRateUpdateMaxSeconds-s.Config.RealRateUpdateMinSeconds) + s.Config.RealRateUpdateMinSeconds
	return time.Second * time.Duration(t)
}

func (s *Service) makeRandomRate(curName, mode string) CurrencyRates {
	buyExchange, sellExchange := s.randExchanges()
	r := CurrencyRates{
		Currency:     curName,
		BuyPrice:     s.randCurrencyPrice(curName, s.Config.BuyMinPer, s.Config.BuyMaxPer),
		SellPrice:    s.randCurrencyPrice(curName, s.Config.SellMinPer, s.Config.SellMaxPer),
		BuyExchange:  buyExchange,
		SellExchange: sellExchange,
		Mode:         mode,
	}
	return r
}

func (s *Service) randCurrencyPrice(currName string, min, max float64) float64 {
	t := s.randGen.Float64()
	t = t*(max-min) + min
	return helper.Round(s.rates[currName]*t, 8)
}

func (s *Service) GetRates(mode string) []CurrencyRates {
	var r []CurrencyRates
	if mode == "real" {
		s.modifiedRealRatesMx.RLock()
		defer s.modifiedRealRatesMx.RUnlock()
		r = make([]CurrencyRates, len(s.modifiedRealRates))
		copy(r, s.modifiedRealRates)
	} else {
		s.modifiedDemoRatesMx.RLock()
		defer s.modifiedDemoRatesMx.RUnlock()
		r = make([]CurrencyRates, len(s.modifiedDemoRates))
		copy(r, s.modifiedDemoRates)
	}
	return r
}

// fetchRates requested all rates from min-api.cryptocompare.com
func (s *Service) fetchRates() error {
	url := "https://min-api.cryptocompare.com/data/pricemulti?tsyms=BTC&fsyms=" + strings.Join(s.currencies, ",")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		logrus.WithError(err).Error("Failed close response body")
	}

	var prices map[string]map[string]float64
	if err = json.Unmarshal(contents, &prices); err != nil {
		return err
	}
	for name, price := range prices {
		if _, ok := price["BTC"]; !ok {
			logrus.Errorf("Unexpected response from cryptocompare.com.\nRaw response: %s", contents)
			continue
		}
		s.rates[name] = price["BTC"]
	}
	return nil
}

func initCurrencies(db *gorm.DB) []string {
	var Currencies []struct {
		Name string
	}
	if err := db.Table("currencies").Find(&Currencies).Error; err != nil {
		panic(err)
	}
	if len(Currencies) < 1 {
		panic("In database lacks rates")
	}
	var resp []string
	for _, e := range Currencies {
		resp = append(resp, e.Name)
	}
	return resp
}

func (s *Service) randExchanges() (string, string) {
	c1 := s.randGen.Intn(len(s.exchanges))
	c2 := s.randGen.Intn(len(s.exchanges))
	if c1 == c2 {
		return s.randExchanges()
	}
	return s.exchanges[c1], s.exchanges[c2]
}

func initExchanges(db *gorm.DB) []string {
	var Exchanges []struct {
		Name string
	}
	if err := db.Table("exchanges").Find(&Exchanges).Error; err != nil {
		panic(err)
	}
	if len(Exchanges) < 2 {
		panic("In database lacks exchanges")
	}
	var response []string
	for _, e := range Exchanges {
		response = append(response, e.Name)
	}
	return response
}
