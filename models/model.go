package models

import (
	"errors"
	"fmt"
	"github.com/caibinsong/wedding/config"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var (
	db       *gorm.DB
	redis_db redis.Conn

	ERROR_DB_ACTION     = errors.New("数据库操作失败")
	ERROR_NOT_FIND_USER = errors.New("用户未注册")
)

func InitDataBase() {
	openDb()
}

func openDb() {
	db1, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Local",
		config.GetConfig().Mysql.User,
		config.GetConfig().Mysql.Pass,
		config.GetConfig().Mysql.Host,
		config.GetConfig().Mysql.DBName,
	))
	if err != nil {
		log.Printf("init DateBase error: [%v]", err)
		return
	}
	if config.GetConfig().Mysql.Debug == "dev" {
		db1.LogMode(true)
	}

	db = db1
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(50)
	initTable()

	//连接redis数据库
	c, err := redis.Dial("tcp", config.GetConfig().Redis)
	if err != nil {
		log.Println("Connect to redis error", err)
		return
	}
	if config.GetConfig().RedisPassWord != "" {
		if _, err := c.Do("AUTH", config.GetConfig().RedisPassWord); err != nil {
			log.Println(err)
		}
	}
	redis_db = c
}

func initTable() {
	db.AutoMigrate(new(RedPacket), new(RedPacketParams), new(Balance), new(BalanceLog), new(Spending), new(Client), new(Recharge))
}

func GetDBObject() *gorm.DB {
	return db
}

func GetRedisDB() redis.Conn {
	return redis_db
}
