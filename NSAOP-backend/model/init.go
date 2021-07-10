package model

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"nsaop/config"
)

var DB *gorm.DB

func Build() {
	DB.AutoMigrate(&User{}, &Refresh{}, &Location{}, &Service{})
}

func Run() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Model.Get("user"),
		config.Model.Get("pwd"),
		config.Model.Get("ip"),
		config.Model.Get("port"),
		config.Model.Get("db"),
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fail to connect to database: %v", err)
	}

	Build()
}
