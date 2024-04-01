package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func NewDatabase(viper *viper.Viper) *sqlx.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD") // viper.GetString("database.password")
	host := os.Getenv("DB_HOST")         //viper.GetString("database.host")
	port := os.Getenv("DB_PORT")         // viper.GetInt("database.port")
	database := os.Getenv("DB_NAME")     //viper.GetString("database.name")
	idleConnection := 100
	maxConnection := 600

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", username, password, host, port, database, os.Getenv("DB_PARAMS"))
	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
	// fmt.Println(dsn)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed connect database: %v", err)
	}
	db.SetMaxIdleConns(idleConnection)
	db.SetMaxOpenConns(maxConnection)

	fmt.Println("Database Connected")

	return db
}
