package main

import (
	"cryptobot/app"
	"cryptobot/conf"
	"fmt"
	"os"
	"time"

	"cryptobot/controller"
	"cryptobot/services/rates"
	"cryptobot/services/robot2"
	"cryptobot/services/trading"

	"context"
	"database/sql"
	"net/http"
	_ "net/http/pprof"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid path to config file.")
		fmt.Println("Please follow next example:")
		fmt.Println("go run main.go ./conf/config.json")
		os.Exit(-1)
	}
	config := conf.FromFile(os.Args[1])

	if config.Debug.Log {
		logrus.SetLevel(logrus.DebugLevel)
	}
	db, err := initDB(config)
	for err != nil {
		logrus.WithError(err).Error("Failed init connection to database. Trying again after ten seconds")
		time.Sleep(time.Second * 10)
		db, err = initDB(config)
	}
	application := app.New(db)

	rateService := rates.New(db, config.Rates)
	tradingService := trading.New(application, rateService, config.Leverages, config.MinBalance, config.FixFactor, config.Rates.DemoRateUpdateMinSeconds, config.Rates.RealRateUpdateMinSeconds)
	robots := robot2.New(tradingService.NewBatchRobotTradeCh, config.Robot, rateService, application)
	c := controller.New(config, rateService, application, tradingService, robots)

	// Run BaseBalanceJob and schedule it every day at 00:10
	baseBalanceJob := BaseBalanceJob{
		db: db,
	}

	realTradesJob := RemoveOldRealTradesJob{
		db: db,
	}

	demoTradesJob := RemoveOldDemoTradesJob{
		db: db,
	}

	scheduler := cron.New()
	scheduler.AddJob("0 10 0 * * *", baseBalanceJob)
	scheduler.AddJob("@daily", realTradesJob)
	scheduler.AddJob("@every 10m", demoTradesJob)
	scheduler.Start()

	go rateService.Start()
	go tradingService.Start()
	go tradingService.BatchTradesHandler()
	go robots.Start()
	go http.ListenAndServe(":8080", nil)
	logrus.Fatal(c.Start())
}

func initDB(conf conf.Main) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", conf.DbConnURL())
	if err != nil {
		return nil, err
	}
	if err := db.DB().Ping(); err != nil {
		return nil, err
	}
	if conf.Debug.DB {
		db = db.Debug()
	}
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}
	if _, err := migrate.Exec(db.DB(), "mysql", migrations, migrate.Up); err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(40)
	return db, nil
}

type BaseBalanceJob struct {
	db *gorm.DB
}

func (j BaseBalanceJob) Run() {
	query := `UPDATE users 
	LEFT JOIN 
		(SELECT user, SUM(profit) AS today_profit FROM trades WHERE created_at >= CURRENT_DATE() GROUP BY user) t 
		ON users.id = t.user 
	SET users.today_base_balance = users.balance - IFNULL(t.today_profit, 0)`

	if err := j.db.Exec(query).Error; err != nil {
		logrus.Errorf("Error when initializing today_base_balance. %v", err)
	}
}

type RemoveOldRealTradesJob struct {
	db *gorm.DB
}

func (j RemoveOldRealTradesJob) Run() {
	query := `DELETE FROM trades WHERE mode = 'real' AND DATE(created_at) <= (CURRENT_DATE() - INTERVAL 3 MONTH)`

	tx, _ := j.db.DB().BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if _, err := tx.Exec(query); err != nil {
		logrus.Errorf("Error when removing old trades for real user. %v", err)
		tx.Rollback()
	}
	tx.Commit()
}

type RemoveOldDemoTradesJob struct {
	db        *gorm.DB
	isRunning bool
}

func (j RemoveOldDemoTradesJob) Run() {
	if !j.isRunning {
		query := `DELETE FROM trades WHERE mode = 'demo' AND created_at <= (NOW() - INTERVAL 12 HOUR) LIMIT 10000000`

		logrus.Infof("Starting demo cleanup")
		j.isRunning = true
		start := time.Now()
		tx, _ := j.db.DB().BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if res, err := tx.Exec(query); err != nil {
			duration := time.Now().Sub(start)
			logrus.Errorf("Error when removing old trades for demo user (%s). %v", duration.String(), err)
			tx.Rollback()
		} else {
			tx.Commit()
			rows, _ := res.RowsAffected()
			duration := time.Now().Sub(start)
			logrus.Infof("Demo cleanup successfull. %d rows deleted on %s", rows, duration)
		}
		j.isRunning = false
	}
}
