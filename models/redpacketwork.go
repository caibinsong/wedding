package models

import (
	"encoding/json"
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

func GetRedPack(userid, rpid int64) error {
	getMutex.Lock()
	defer getMutex.Unlock()

	return work.getRedPack(userid, rpid)
}

func (this *redPacketWork) getRedPack(user_id, rp_id int64) error {
	status := CheckRedPacket(user_id, rp_id)
	if status == 0 {
		redPacket, err := redis.String(GetRedisDB().Do("GET", fmt.Sprintf("%s%d", config.REDIS_REDPACK, rp_id)))
		if err != nil {
			if strings.Index(err.Error(), "connection timed out") > 0 {
				ConnectRedis()
				redPacket, err = redis.String(GetRedisDB().Do("GET", fmt.Sprintf("%s%d", config.REDIS_REDPACK, rp_id)))
			}
			if err != nil {
				log.Println("redis get:", err)
				return err
			}
			log.Println(redPacket, err)
			return err
		}
		redPacketUser, err := redis.String(GetRedisDB().Do("GET", fmt.Sprintf("%s%d", config.REDIS_REDPACK_USER, rp_id)))
		if err != nil {
			return err
		}
		//1_3.59;2_7.93;3_3.48;
		index := strings.Index(redPacket, ";")
		if index < 0 {
			return fmt.Errorf(statusMap[1])
		}
		sRedPacket := redPacket[:index]
		if strings.Index(sRedPacket, "_") > 0 {
			arr := strings.Split(sRedPacket, "_")
			if len(arr) == 2 {
				id, _ := strconv.Atoi(arr[0])
				money, _ := strconv.ParseFloat(arr[1], 64)
				err = GrabRedPacket(user_id, rp_id, int64(id), money)
				if err != nil {
					log.Println(err.Error())
					return fmt.Errorf("红包异常")
				}
				_, err = GetRedisDB().Do("SET", fmt.Sprintf("%s%d", config.REDIS_REDPACK, rp_id), redPacket[index+1:])
				if err != nil {
					log.Println("redis set failed:", err)
					return err
				}
				_, err = GetRedisDB().Do("SET", fmt.Sprintf("%s%d", config.REDIS_REDPACK_USER, rp_id),
					fmt.Sprintf("%s;%d", redPacketUser, user_id))
				if err != nil {
					log.Println("redis set failed:", err)
					return err
				}

				return nil
			} else {
				log.Println("getRedPack err:", sRedPacket)
				return fmt.Errorf("红包异常")
			}
		} else {
			log.Println("getRedPack err:", sRedPacket)
			return fmt.Errorf("红包异常")
		}
	}
	return fmt.Errorf(statusMap[status])
}

type AccessCtr struct {
	RpId   int64
	UserId int64
}

type AccessCtrlWork struct {
	ServerName string
	Methodname string
	Data       chan AccessCtr // map[string]interface{}
}

func StartAccessCtrWork(serverName, methodname string) {
	accessCtrWork.ServerName = serverName
	accessCtrWork.Methodname = methodname
	accessCtrWork.Data = make(chan AccessCtr /*map[string]interface{}*/, 5000)
	for i := 0; i < 50; i++ {
		go accessCtrWork.work(i)
	}
}

func AddAccessCtrWork(data AccessCtr /*map[string]interface{}*/) {
	accessCtrWork.Data <- data
}

func (this *AccessCtrlWork) work(i int) {
	defer func() {
		log.Println("work", i, "is die")
	}()
	log.Println("AccessCtr", i, "work start")
	httpClient := utils.NewHttpClient()
	Msg := ""
	for {
		accessCtr := <-this.Data
		Msg = ""

		redPacket, err := FindRedPacketByRpId(accessCtr.RpId)
		if err != nil {
			log.Println(err.Error())
			Msg = "红包已经抢完！"
		}

		speeding, err := QuerySpending(redPacket.UserId, redPacket.CreateAt)
		if err != nil {
			log.Println(err.Error())
			Msg = "数据库操作失败"
		}
		// 组织广播内容
		data := map[string]interface{}{"rp_id": accessCtr.RpId,
			"red_type": redPacket.RedPacketType}
		var msg map[string]interface{}
		if Msg == "" {
			msg = map[string]interface{}{"code": 0}
		} else {
			msg = map[string]interface{}{"code": 1, "msg": Msg}
		}
		roomSvr := map[string]interface{}{"chatroomId": redPacket.RoomId,
			"weddingId": speeding.WeddingId,
			"userId":    accessCtr.UserId,
			"data":      data,
			"msg":       msg}

		bRoomSvr, err := json.Marshal(map[string]interface{}{"msg": roomSvr})
		if err != nil {
			log.Println(err.Error())
		}
		content := map[string]interface{}{
			"type":    2,
			"content": string(bRoomSvr),
		}
		bContent, err := json.Marshal(content)
		if err != nil {
			log.Println(err.Error())
		}
		body := map[string]interface{}{
			"type":    "HLBUser",
			"idList":  []int64{accessCtr.UserId},
			"content": string(bContent),
		}

		err = httpClient.AccessCtrlSvr(this.ServerName, this.Methodname, body)
		if err != nil {
			log.Println(err)
		}
	}
}
