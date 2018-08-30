package main

import (
	"cryptobot/conf"
	"database/sql"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rubenv/sql-migrate"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		printHelpAndExit()
	}
	c := conf.FromFile("./conf/config.json")

	db, err := sql.Open("mysql", c.DbConnURL())
	if err != nil {
		panic(err)
	}

	if os.Args[1] == "current" {
		fmt.Println("\n\nLast active migration:", currentMigration(db))
		return
	}

	var n int
	if len(os.Args) == 3 {
		n, err = strconv.Atoi(os.Args[2])
		if err != nil || n < 1 {
			printHelpAndExit()
		}
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}
	if n == 0 {
		switch os.Args[1] {
		case "up":
			_, err = migrate.Exec(db, "mysql", migrations, migrate.Up)
		case "down":
			_, err = migrate.Exec(db, "mysql", migrations, migrate.Down)
		default:
			printHelpAndExit()
		}
	} else {
		switch os.Args[1] {
		case "up":
			_, err = migrate.ExecMax(db, "mysql", migrations, migrate.Up, n)
		case "down":
			_, err = migrate.ExecMax(db, "mysql", migrations, migrate.Down, n)
		default:
			printHelpAndExit()
		}
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("\n\nMigrations successfully applied! Last active migration:", currentMigration(db))
}

func printHelpAndExit() {
	fmt.Println("\n\nInvalid input.")
	fmt.Println("Parameter <type> required. Must be up/down/current")
	fmt.Println("<current> returns current last active migrate name.")
	fmt.Println("Parameter <n> optional. Determines, how many steps upping/downing.")
	fmt.Println("If <n> not selected - go to top/bottom migrate.")
	fmt.Println("Example input: make migrate type=up n=1")
	os.Exit(-1)
}

func currentMigration(db *sql.DB) string {
	row := db.QueryRow("SELECT id from gorp_migrations order by id desc limit 1")
	var m string
	if err := row.Scan(&m); err != nil {
		return "Not found appied migrations!\n"
	}
	return m
}
