package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Simple function to get a connection to query the db
func GetDb() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env error")
	}

	dbCfg := mysql.Config{
		User:                 os.Getenv("DbUser"),
		Passwd:               os.Getenv("DbPassword"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DbAddr"),
		DBName:               os.Getenv("DbName"),
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		log.Fatal(err.Error())
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	return db
}
