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
	rp, err := FindRedPacketInfoByRpId(rpId)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = UpdateRecharge(rp.UserId, rp.CreateAt, Transaction_id)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = db.Table("cRedPacketParams").Where(&RedPacketParams{RpId: rpId}).Updates(map[string]interface{}{"status": 1}).Error
	if err != nil {
		log.Println(err.Error())
	}
	err = db.Table("cRedPacket").Where(&RedPacket{RpId: rpId}).Updates(map[string]interface{}{"status": 1}).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return err
}

//通过红包ID查询红包信息
func FindRedPacketInfoByRpId(rpId int64) (RedPacket, error) {
	redPacket := RedPacket{}
	err := db.Where(&RedPacket{RpId: rpId}).First(&redPacket).Error
	if err != nil {
		log.Println(err.Error())
	}
	return redPacket, err
}

//通过红包ID查询红包明细信息
func FindRedPacketParamsByRpId(rpId int64) ([]RedPacketParams, error) {
	redPacketParamsList := make([]RedPacketParams, 0)
	err := db.Where(fmt.Sprintf("cRedPacketParams.rp_id=%d and cRedPacketParams.user_id<>0", rpId)).
		Find(&redPacketParamsList).Error

	if err != nil {
		log.Println(err.Error())
	}
	return redPacketParamsList, nil
}

func CheckUserRedPacket(userId int64, req *config.Req_RedPacket) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	rp, err := FindRedPacketInfoByRpId(req.Data.RpId)
	if err != nil {
		return nil, err
	}
	if rp.Status == 0 || time.Now().Unix() > (rp.CreateAt+config.ONE_DAY) {
		resultMap["status"] = 3
		return resultMap, nil
	}
	resultMap["status"] = CheckRedPacket(userId, req.Data.RpId)
	return resultMap, nil
}

func CheckRedPacket(user_id, rp_id int64) int {
	if _, ok := isEndStatus(rp_id); ok {
		return 3
	}
	if _, ok := hasRedPacket(rp_id); !ok {
		return 1
	}
	if _, ok := repeatUser(user_id, rp_id); ok {
		return 2
	}
	return 0
}

//如果redis中 空就是 过期了
func isEndStatus(rp_id int64) (string, bool) {
	redPacket, err := redis.String(GetRedisDB().Do("GET",
		fmt.Sprintf("%s%d", config.REDIS_REDPACK, rp_id)))
	if err != nil {
		return redPacket, true
	}
	return redPacket, false
}

//判断红包是否还有
func hasRedPacket(rp_id int64) (string, bool) {
	redPacket, err := redis.String(GetRedisDB().Do("GET",
		fmt.Sprintf("%s%d", config.REDIS_REDPACK, rp_id)))
	if err != nil {
		return "", false
	}
	if strings.Index(redPacket, "_") > 0 {
		return "", true
	}
	return "", false
}

//判断红包是否已经抢过了
func repeatUser(userid, rp_id int64) (string, bool) {
	getRedPacketUser, err := redis.String(GetRedisDB().Do("GET", fmt.Sprintf("%s%d", config.REDIS_REDPACK_USER, rp_id)))
	if err != nil {
		return "", false
	}
	return getRedPacketUser, in(getRedPacketUser, strconv.Itoa(int(userid)))
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

func FindRedPacketByRpId(rpId int64) RedPacket {
	redPacket := RedPacket{}
	db.Where(&RedPacket{RpId: rpId}).First(&redPacket)
	return redPacket
}
