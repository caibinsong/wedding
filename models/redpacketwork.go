package models

import (
	"fmt"
	"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/utils"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"strings"
	"sync"
)

var (
	work          *redPacketWork = &redPacketWork{}
	getMutex      sync.Mutex
	statusMap     map[int]string  = map[int]string{1: "红包已抢完", 2: "已经抢过红包了", 3: "红包已过期"}
	accessCtrWork *AccessCtrlWork = &AccessCtrlWork{}
)

type redPacketWork struct {
}

func GetRedPack(userid, rpid int64) (int, float64, error) {
	getMutex.Lock()
	defer getMutex.Unlock()
	return work.getRedPack(userid, rpid)
}

func (this *redPacketWork) getRedPack(user_id, rp_id int64) (int, float64, error) {
	status := CheckRedPacket(user_id, rp_id)
	if status == 0 {
		redPacket, err := redis.String(GetRedisDB().Do("GET",
			fmt.Sprintf("%s%d", config.REDIS_REDPACK, rp_id)))
		if err != nil {
			return 0, 0, err
		}
		redPacketUser, err := redis.String(GetRedisDB().Do("GET",
			fmt.Sprintf("%s%d", config.REDIS_REDPACK_USER, rp_id)))
		if err != nil {
			return 0, 0, err
		}
		//1_3.59;2_7.93;3_3.48;
		index := strings.Index(redPacket, ";")
		if index < 0 {
			return 0, 0, fmt.Errorf(statusMap[1])
		}
		sRedPacket := redPacket[:index]
		if strings.Index(sRedPacket, "_") > 0 {
			arr := strings.Split(sRedPacket, "_")
			if len(arr) == 2 {
				id, _ := strconv.Atoi(arr[0])
				money, _ := strconv.ParseFloat(arr[1], 64)

				_, err = GetRedisDB().Do("SET", fmt.Sprintf("%s%d", config.REDIS_REDPACK, rp_id),
					redPacket[index+1:], "EX", "86400")
				if err != nil {
					log.Println("redis set failed:", err)
					return 0, 0, err
				}

				_, err = GetRedisDB().Do("SET", fmt.Sprintf("%s%d", config.REDIS_REDPACK_USER, rp_id),
					fmt.Sprintf("%s;%d", redPacketUser, user_id), "EX", "86400")
				if err != nil {
					log.Println("redis set failed:", err)
					return 0, 0, err
				}

				return id, money, nil
			} else {
				log.Println("getRedPack err:", sRedPacket)
				return 0, 0, fmt.Errorf("红包异常")
			}
		} else {
			log.Println("getRedPack err:", sRedPacket)
			return 0, 0, fmt.Errorf("红包异常")
		}
	}
	return 0, 0, fmt.Errorf(statusMap[status])
}

type AccessCtrlWork struct {
	ServerName string
	Methodname string
	Data       chan map[string]interface{}
}

func StartAccessCtrWork(serverName, methodname string) {
	accessCtrWork.ServerName = serverName
	accessCtrWork.Methodname = methodname
	accessCtrWork.Data = make(chan map[string]interface{})
	//var wg *sync.WaitGroup = &sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		go accessCtrWork.work(i)
		//wg.Add(1)
	}
	//wg.Wait()
}

func AddAccessCtrWork(data map[string]interface{}) {
	accessCtrWork.Data <- data
}

func (this *AccessCtrlWork) work(i int) {
	defer func() {
		log.Println("work", i, "is die")
	}()
	for {
		data := <-this.Data
		err := utils.NewHttpClient().AccessCtrlSvr(this.ServerName, this.Methodname, data)
		if err != nil {
			log.Println(err)
		}
	}
}
