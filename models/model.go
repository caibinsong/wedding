package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/Amniversary/wedding-logic-redpacket/config"
	"log"
)

var db *gorm.DB

func InitDataBase() {
	openDb()
}

func openDb() {
	db1, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Local",
		config.USER,
		config.PASS,
		config.HOST,
		config.DBName,
	))
	if err != nil {
		log.Printf("init DateBase error: [%v]", err)
		return
	}
	if config.DEBUG == "dev" {
		db1.LogMode(true)
	}

	db = db1
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(50)
	initTable()
}

func initTable() {
	db.AutoMigrate(new(RedPacket), new(RedPacketParams))
}
