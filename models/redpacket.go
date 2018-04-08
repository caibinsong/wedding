package models

import (
	"fmt"
	"github.com/caibinsong/wedding/config"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

type RedPacket struct {
	RpId           int64   `gorm:"primary_key" json:"rp_id"`
	Guid           string  `gorm:"not null; default:''; type:varchar(128)" json:"guid"`              //GUID
	UserId         int64   `gorm:"not null; default:0; type:int; index" json:"user_id"`              //用户ID
	RoomId         int64   `gorm:"not null; default:0; type:int; index" json:"room_id"`              //房间ID
	RedPacketNum   int64   `gorm:"not null; default:0; type:int" json:"red_packet_num"`              //红包数量
	RedPacketMoney float64 `gorm:"not null; default:0; type:decimal(12,2)" json:"red_packet_money"`  //红包金额
	RedPacketType  int64   `gorm:"not null; default:0; type:int" json:"red_packet_type"`             //红包类型
	GetNum         int64   `gorm:"not null; default:0; type:int" json:"get_num"`                     //获取的红包数
	Status         int64   `gorm:"not null; default:0; type:int" json:"status"`                      //红包状态
	EndStatus      int64   `gorm:"not null; default:0; type:int" json:"end_status"`                  //过期时间
	CreateAt       int64   `gorm:"not null; default:0; type:int; index" json:"create_at"`            //问题红包
	Question       string  `gorm:"not null; type:text" json:"question"`                              //创建时间
	Remark1        string  `gorm:"not null; default:''; type:varchar(128); column:remark1" json:"-"` //备注1
	Remark2        string  `gorm:"not null; default:''; type:varchar(128); column:remark2" json:"-"` //备注2
}

type RedPacketParams struct {
	ID             int64   `gorm:"primary_key" json:"id"`
	RpId           int64   `gorm:"not null; default:0; type:int; index" json:"rp_id"`                //红包ID
	RedPacketNo    int64   `gorm:"not null; default:0; type:int; index" json:"red_packet_no"`        //红包序号
	RedPacketMoney float64 `gorm:"not null; default:0; type:decimal(12,2)" json:"red_packet_money"`  //红包金额
	Status         int64   `gorm:"not null; default:0; type:int" json:"status"`                      //红包状态
	Lucky          int64   `gorm:"not null; default:0; type:int" json:"lucky"`                       //最佳手气
	UserId         int64   `gorm:"not null; default:0; type:int" json:"user_id"`                     //用户ID
	UpdateAt       int64   `gorm:"not null; default:0; type:int; index" json:"update_at"`            //领取时间
	Remark1        string  `gorm:"not null; default:''; type:varchar(128); column:remark1" json:"-"` //备注1
	Remark2        string  `gorm:"not null; default:''; type:varchar(128); column:remark2" json:"-"` //备注2
}

const (
	SQL_FindRedPacketByRpId       = "SELECT * FROM `cRedPacket`  WHERE `rp_id` = %d"
	SQL_FindRedPacketParamsByRpId = "SELECT * FROM `cRedPacketParams`  WHERE (cRedPacketParams.rp_id=%d and cRedPacketParams.user_id<>0)"
	SQL_UpRedPacketParams_Status  = "UPDATE `cRedPacketParams` SET `status` = 1  WHERE `rp_id` = %d"
	SQL_UpRedPacket_Status        = "UPDATE `cRedPacket` SET `status` = 1  WHERE `rp_id` = %d"
	//抢红包
	GRAB_RED_PACKET        = "update  cRedPacket  set get_num=get_num+1  where rp_id=%d and `status`=1;"
	GRAB_RED_PACKET_PARAMS = "update cRedPacketParams set user_id=%d, update_at=%d  where rp_id=%d and red_packet_no=%d and `status`=1"
)

func (RedPacket) TableName() string {
	return "cRedPacket"
}

func (RedPacketParams) TableName() string {
	return "cRedPacketParams"
}

//微信回调时 修改redpacket 和redpacketparams 的status状态
func UpDateRedPacketStatus(rpId int64, Transaction_id string) error {
	rp, err := FindRedPacketByRpId(rpId)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = UpdateRecharge(rp.UserId, rp.CreateAt, Transaction_id)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = db.Exec(fmt.Sprintf(SQL_UpRedPacketParams_Status, rpId)).Error
	if err != nil {
		log.Println(err.Error())
	}
	err = db.Exec(fmt.Sprintf(SQL_UpRedPacket_Status, rpId)).Error
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

//通过红包ID查询红包信息
func FindRedPacketByRpId(rpId int64) (RedPacket, error) {
	redPacket := RedPacket{}
	err := db.Raw(fmt.Sprintf(SQL_FindRedPacketByRpId, rpId)).Scan(&redPacket).Error
	if err != nil {
		log.Println(err.Error())
	}
	return redPacket, err
}

//通过红包ID查询红包明细信息
func FindRedPacketParamsByRpId(rpId int64) ([]RedPacketParams, error) {
	redPacketParamsList := make([]RedPacketParams, 0)
	rows, err := db.Raw(fmt.Sprintf(SQL_FindRedPacketParamsByRpId, rpId)).Rows()
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	} else {
		for rows.Next() {
			redPacketParams := RedPacketParams{}
			rows.Scan(&redPacketParams.ID, &redPacketParams.RpId, &redPacketParams.RedPacketNo,
				&redPacketParams.RedPacketMoney, &redPacketParams.Status, &redPacketParams.Lucky,
				&redPacketParams.UserId, &redPacketParams.UpdateAt, &redPacketParams.Remark1, &redPacketParams.Remark2)

			redPacketParamsList = append(redPacketParamsList, redPacketParams)
		}
	}
	return redPacketParamsList, nil
}

func CheckUserRedPacket(userId int64, req *config.Req_RedPacket) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})

	resultMap["status"] = CheckRedPacket(userId, req.Data.RpId)
	return resultMap, nil
}

func CheckRedPacket(user_id, rp_id int64) int {
	if isEndStatus(rp_id) {
		return 3
	}
	if !hasRedPacket(rp_id) {
		return 1
	}
	if repeatUser(user_id, rp_id) {
		return 2
	}
	return 0
}

//如果redis中 空就是 过期了
func isEndStatus(rp_id int64) bool {
	rp, err := FindRedPacketByRpId(rp_id)
	if err != nil {
		return true
	}
	if rp.Status == 0 || time.Now().Unix() > (rp.CreateAt+config.ONE_DAY) {
		return true
	}
	return false
}

//判断红包是否还有
func hasRedPacket(rp_id int64) bool {
	redPacket, err := redis.String(GetRedisDB().Do("GET", fmt.Sprintf("%s%d", config.REDIS_REDPACK, rp_id)))
	if err != nil {
		if strings.Index(err.Error(), "connection timed out") > 0 {
			ConnectRedis()
			redPacket, err = redis.String(GetRedisDB().Do("GET", fmt.Sprintf("%s%d", config.REDIS_REDPACK, rp_id)))
		}
		if err != nil {
			return false
		}
	}
	if strings.Index(redPacket, "_") > 0 {
		return true
	}
	return false
}

//判断红包是否已经抢过了
func repeatUser(userid, rp_id int64) bool {
	getRedPacketUser, err := redis.String(GetRedisDB().Do("GET", fmt.Sprintf("%s%d", config.REDIS_REDPACK_USER, rp_id)))
	if err != nil {
		return false
	}
	return in(getRedPacketUser, strconv.Itoa(int(userid)))
}

func in(all, one string) bool {
	arr := strings.Split(all, ";")
	for _, v := range arr {
		if v == one {
			return true
		}
	}
	return false
}
