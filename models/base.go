package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	envPath := filepath.Join("utils", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Print(err)
	}
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, dbHost, dbPort, dbName)
	fmt.Println("databse URI ", dbURI)
	conn, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	db = conn
	db.Debug().AutoMigrate(&Contact{}, &Account{})
	fmt.Println("connection was successful")
}

func GetDB() *gorm.DB {
	return db
}
