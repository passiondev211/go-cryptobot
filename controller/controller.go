package controller

import (
	"bytes"
	"cryptobot/app"
	"cryptobot/conf"
	"cryptobot/services/rates"
	"cryptobot/services/robot2"
	"cryptobot/services/trading"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/orange-cloudfoundry/ipfiltering"
)

// Controller implemented router of project
type Controller struct {
	Config conf.Main
	App    *app.App
	Gin    *gin.Engine

	// Service which executed all incoming from robots trades
	Core *trading.Service

	// Associated map userID -> his trading robot
	Robots *robot2.Service

	// Rates stories info about current currencies
	Rates *rates.Service

	// Dashboard includes all gathering info about user
	Dashboard Dashboard

	// IP Filter
	ipFilter *ipfiltering.IpFiltering
}

// MyCustomClaims includes fields from jwt token.
type MyCustomClaims struct {
	UserID   int64  `json:"id"`
	Email    string `json:"email"`
	Language string `json:"language"`
	jwt.StandardClaims
}

func maxRequestsAtOnce(n int) gin.HandlerFunc {
	s := make(chan struct{}, n)
	return func(c *gin.Context) {
		s <- struct{}{}
		defer func() { <-s }()
		c.Next()
	}
}

var allowedMethods string = "GET, PUT, POST, DELETE"

func initGin() *gin.Engine {
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(maxRequestsAtOnce(50))
	g.Use(cors.Middleware(cors.Config{
		Origins:        "http://localhost:3000, https://profitcoins.io, https://app.profitcoins.io, https://api.profitcoins.io, https://dev.profitcoins.io, https://app.dev.profitcoins.io, https://api.dev.profitcoins.io",
		Methods:        allowedMethods,
		RequestHeaders: "Origin, Authorization, Content-Type, Access-Control-Allow-Origin",
		ExposedHeaders: "",
		MaxAge:         50 * time.Second,
	}))

	return g
}

func security(ctx *gin.Context) {
	if !strings.Contains(allowedMethods, ctx.Request.Method) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Method %s is not allowed", ctx.Request.Method),
		})
		ctx.Abort()
	}
	if ctx.Request.Method != "GET" && ctx.ContentType() != "application/json" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Content Type %s is not allowed", ctx.ContentType()),
		})
		ctx.Abort()
	}
	ctx.Next()
}

func New(config conf.Main, rateService *rates.Service, application *app.App, tradingService *trading.Service, robots *robot2.Service) *Controller {
	filter := ipfiltering.New(ipfiltering.Options{
		BlockByDefault: true,
		AllowedIPs:     config.Security.AllowedAdminIP,
	})
	return &Controller{
		Config:   config,
		Gin:      initGin(),
		Rates:    rateService,
		App:      application,
		Core:     tradingService,
		Robots:   robots,
		ipFilter: filter,
	}
}

// Start initialize endpoints and launch http server.
func (c *Controller) Start() error {

	c.Gin.Use(c.logHttp)

	c.Gin.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	c.Gin.GET("/api/v1/login", c.login)

	// CMS endpoints. For access to this endpoints need access_token in header
	admin := c.Gin.Group("")
	admin.Use(security)
	admin.Use(c.filterIPAddress)
	admin.Use(c.checkAdminToken)
	admin.POST("/api/v1/signup", c.signUp)
	admin.POST("/api/v1/signin", c.signIn)
	admin.POST("/api/v1/deposit", c.deposit)
	admin.POST("/api/v1/withdraw", c.withdraw)
	admin.GET("/api/v1/users", c.getUsers)
	admin.GET("/api/v1/users/:userID", c.getUser)
	admin.PUT("/api/v1/users/:userID", c.changeMargin)
	admin.PUT("/api/v1/users/:userID/trade-bounds", c.changeTradeBounds)
	admin.DELETE("/api/v1/users/:userID/trade-bounds", c.resetTradeBounds)
	admin.PUT("/api/v1/users/:userID/fix-factor", c.changeFixFactor)
	admin.DELETE("/api/v1/users/:userID/fix-factor", c.resetFixFactor)
	admin.GET("/api/v1/config", c.getConfig)
	admin.POST("/api/v1/config", c.setConfig)

	// We can't use just a.Gin.Serve here because Gin will say / collides with /v1 stuff.
	c.Gin.Use(c.checkCookie)
	c.Gin.Use(static.Serve("/", static.LocalFile("./dist", true)))
	c.Gin.Use(static.Serve("/trades", static.LocalFile("./dist", true)))

	private := c.Gin.Group("")
	private.Use(security)
	private.Use(c.authentication)
	private.GET("/api/v1/exchange-rates", c.exchangeRates)
	private.GET("/api/v1/user-info", c.userInfo)
	private.GET("/api/v1/trade-stats", c.tradeStats)
	private.GET("/api/v1/trades", c.retrieveTrades)
	private.POST("/api/v1/trades", c.createTrade)
	private.GET("/api/v1/app-config", c.appConfig)
	private.POST("/api/v1/bot-toggle", c.startTrading)
	private.GET("/api/v1/dashboard-info", c.dashboardInfo)
	private.GET("/api/v1/tour-visited", c.getTourVisited)
	private.POST("/api/v1/tour-visited", c.setTourVisited)
	private.GET("/api/v1/notification", c.getNotification)
	private.GET("/api/v1/gtm-events", c.getGTMEvent)
	private.GET("/api/v1/chat-id", c.getChatID)
	private.POST("/api/v1/chat-id", c.setChatID)
	private.POST("/api/v1/logout", c.logout)
	return c.Gin.Run(":" + c.Config.Port)
}

func (c *Controller) authentication(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if tokenString == "" {
		ctx.Redirect(http.StatusFound, c.Config.FE.IJRP)
		ctx.Abort()
		return
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.Config.Auth.JWTSecret), nil
	})
	if err != nil {
		ctx.Redirect(http.StatusFound, c.Config.FE.IJRP)
		ctx.Abort()
		return
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		u := app.User{}
		if err := c.App.DB.Where("id = ?", claims.UserID).First(&u).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Set("user", &u)
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "can't get fields from jwt token",
		})
		ctx.Abort()
		return
	}
	ctx.Next()
}

func (c *Controller) filterIPAddress(ctx *gin.Context) {
	clientIP := ctx.ClientIP()
	if c.ipFilter.Blocked(clientIP) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("ip %s is not allowed", clientIP),
		})
		ctx.Abort()
	}
	ctx.Next()
}

func (c *Controller) checkAdminToken(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Authorization")
	if accessToken != "Bearer "+c.Config.Auth.AccessToken {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "access_token is not valid",
		})
		ctx.Abort()
	}
}

func (c *Controller) checkCookie(ctx *gin.Context) {
	if ctx.Request.URL.String() != "/" &&
		ctx.Request.URL.String() != "/#/" &&
		ctx.Request.URL.String() != "/trades" {
		return
	}

	var sToken string
	for _, cookie := range ctx.Request.Cookies() {
		if cookie.Name == "session_token" {
			sToken = cookie.Value
		}
	}
	if sToken == "" {
		ctx.Redirect(http.StatusFound, c.Config.FE.IJRP)
		ctx.Abort()
	}
	if _, err := jwt.ParseWithClaims(sToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.Config.Auth.JWTSecret), nil
	}); err != nil {
		ctx.Redirect(http.StatusFound, c.Config.FE.IJRP)
		ctx.Abort()
	}
}

func userFromContext(ctx *gin.Context) (*app.User, bool) {
	iu, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "not found user id in context",
		})
		return nil, false
	}
	u, ok := iu.(*app.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "not found user id in context",
		})
		return nil, false
	}
	return u, true
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (c *Controller) logHttp(ctx *gin.Context) {
	// Only log Login or Non GET request
	if strings.HasPrefix(ctx.Request.URL.String(), "/api/v1/login") || (ctx.Request.Method != "GET" && ctx.Request.Method != "") {
		start := time.Now()
		req, _ := httputil.DumpRequest(ctx.Request, true)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		ctx.Next()
		resp := blw.body.String()
		respTime := time.Now().Sub(start)
		log := app.HttpLog{
			ClientIP:     ctx.ClientIP(),
			URL:          ctx.Request.URL.String(),
			Request:      string(req),
			Response:     resp,
			CreatedAt:    time.Now(),
			ResponseTime: respTime.String(),
		}

		go func() {
			if err := c.App.DB.Save(&log).Error; err != nil {
				// Do nothing
			}
		}()
	}
}
