package database

import (
	"fmt"
	m "go-be/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
	RDConn *redis.Client
)

func DB_Init() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		"root",
		"",
		"127.0.0.1",
		"3306",
		"golang_prj",
	)

	var err error
	DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Success : Success to connect database.")

	DBConn.AutoMigrate(&m.User{})
}

func RD_Init() {
	RDConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})
}
