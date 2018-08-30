package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Main include info about API settings.
type Main struct {
	Port           string  `json:"port"`
	MaxTradesCount int     `json:"maxTradesPerDay"`
	MaxMargin      float64 `json:"maxMarginPerDay"`

	Auth        Auth           `json:""`
	Rates       Rates          `json:"rates"`
	DB          DB             `json:"database"`
	FE          AppConfig      `json:"frontEnd"`
	DemoBalance float64        `json:"demoBalance"`
	MinBalance  float64        `json:"minBalance"`
	Leverages   []LeverageConf `json:"leverages"`
	Debug       Debug          `json:"debug"`
	Robot       Robot          `json:"robot"`
	FixFactor   FixFactor      `json:"fixFactor"`
	Security Security `json:"security"`
}

// Auth includes all authorization settings
type Auth struct {
	AccessToken  string `json:"accessToken"`         // CMS token
	AuthTokenTTL int    `json:"authTokenTtlSeconds"` // Time to life of user authentication link, seconds
	JWTSecret    string `json:"jwtSecret"`           // Secret key for generating user token
	JwtTtl       int    `json:"jwtTokenTtlMinutes"`  // Time to life of user session token
}

// Robot configs
type Robot struct {
	TimeoutAfterNewTrade int `json:"timeoutAfterNewTradeS"`  // Timeout between sending new trades, seconds
	ExecTime             int `json:"timeToExecutingTradeMS"` // Time to creating and sending new trade to core, milliseconds
}

// Rates include all rate-service settings
type Rates struct {
	BuyMinPer  float64 `json:"buyingMinPercents"`
	BuyMaxPer  float64 `json:"buyingMaxPercents"`
	SellMinPer float64 `json:"sellingMinPercents"`
	SellMaxPer float64 `json:"sellingMaxPercents"`

	RatesUpdate int `json:"tableUpdateIntervalSeconds"`

	DemoRateUpdateMinSeconds int `json:"demoRatesUpdateIntervalMinSeconds"`
	DemoRateUpdateMaxSeconds int `json:"demoRatesUpdateIntervalMaxSeconds"`

	RealRateUpdateMinSeconds int `json:"realRatesUpdateIntervalMinSeconds"`
	RealRateUpdateMaxSeconds int `json:"realRatesUpdateIntervalMaxSeconds"`
}

// LeverageConf includes info about leverage
type LeverageConf struct {
	RequiredBalance float64 `json:"minBalance"`
	Leverage        int     `json:"leverage"`
	DailyProfit     float64 `json:"profitPerDay"`

	DailyTradeBoundLower  float64 `json:"dailyTradeBoundLower"`
	DailyTradeBoundUpper  float64 `json:"dailyTradeBoundUpper"`
	ResultTradeBoundLower float64 `json:"resultTradeBoundLower"`
	ResultTradeBoundUpper float64 `json:"resultTradeBoundUpper"`
}

// DB include info about database connection.
type DB struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

// Debug include all debugging flags
type Debug struct {
	Log bool `json:"log"`
	DB  bool `json:"db"`
}

// AppConfig includes info, required for front-end
type AppConfig struct {
	APITimeout int `json:"apiTimeout"` // Timeout for front-end requests

	// FIXME(tiabc): Abbreviations don't make the code more readable.
	LRP  string `json:"logoutRedirectPath"`
	DRP  string `json:"depositRedirectPath"`
	FLRP string `json:"firstLinkRedirectPath"`
	SLRP string `json:"secondLinkRedirectPath"`
	FLN  string `json:"firstLinkName"`
	SLN  string `json:"secondLinkName"`
	IJRP string `json:"invalidJwtRedirectPath"`
}

type FixFactor struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
}

type Security struct {
	AllowedAdminIP []string `json:"allowed_admin_ip"`
}

// FromFile read config from path and create configuration struct.
func FromFile(path string) Main {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic("error reading config " + path + ": " + err.Error())
	}
	var conf Main
	if err := json.Unmarshal(bytes, &conf); err != nil {
		panic("error parsing config " + path + ": " + err.Error())
	}
	if len(conf.Leverages) == 0 {
		panic("not found leverages config, please use ./conf/config.sample.json as example")
	}
	return conf
}

func Save(path string, config Main) error {
	configJson, _ := json.MarshalIndent(config, "", "  ")
	if err := ioutil.WriteFile(path, configJson, 0644); err != nil {
		return err
	}

	return nil
}

// DbConnURL create url to connect with database.
func (c *Main) DbConnURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name)
}
