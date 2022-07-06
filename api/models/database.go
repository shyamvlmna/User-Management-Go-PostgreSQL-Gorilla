package models

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDb() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	dbName := os.Getenv("dbName")
	dbPassword := os.Getenv("dbPassword")
	dbUser := os.Getenv("dbUser")
	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)

	Db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Db opened")
	}

	return Db
}
