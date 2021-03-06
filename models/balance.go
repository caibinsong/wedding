package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/utils"
	"github.com/jinzhu/gorm"
	"log"
	"strings"
	"time"
)

const (
	SIGN = `MD5(CONCAT(
                    'id=', id,
                    '&user_id=', user_id,
                    '&balance=', balance,
                    '&recharge_num=', recharge_num,
                    '&reward_num=', reward_num,
                    '&withdraw_cash=', withdraw_cash,
                    '&rand_str=', rand_str,
                    '&chise1bht0z=lkc12i8xzh4wnmz90qnmxca2zqwdc9wqxxzjstlq'))`
)

type Balance struct {
	Id           int64   `gorm:"primary_key" json:"id"`
	UserId       int64   `gorm:"not null; default:0;  type:int" json:"user_id"`                       //用户ID
	Balance      float64 `gorm:"not null; default:0;  type:decimal(12,2)" index" json:"balance"`      //余额
	RechargeNum  float64 `gorm:"not null; default:0;  type:decimal(12,2)" index" json:"recharge_num"` //累计收入
	RewardNum    float64 `gorm:"not null; default:0;  type:decimal(12,2)" json:"reward_num"`          //累计支出
	WithdrawCash float64 `gorm:"not null; default:0;  type:decimal(12,2)" json:"withdraw_cash"`       //累计提现
	Status       int64   `gorm:"not null; default:0;  type:int" json:"status"`                        //账户状态
	RandStr      string  `gorm:"not null; default:'';  type:varchar(100); json:"rand_str"`            //随机串
	Sign         string  `gorm:"not null; default:'';  type:varchar(100); json:"sign"`                //签名串
}

type BalanceLog struct {
	LogId         int64   `gorm:"primary_key" json:"log_id"`
	UserId        int64   `gorm:"not null; default:0;   type:int" json:"user_id"`                   //用户ID
	BalanceId     int64   `gorm:"not null; default:0;   type:int" index" json:"balance_id"`         //账户ID
	OperateType   int64   `gorm:"not null; default:0;   type:int" index" json:"operate_type"`       //操作类型
	OperateValue  float64 `gorm:"not null; default:0;   type:decimal(12,2)" json:"operate_value"`   //操作余额
	BeforeBalance float64 `gorm:"not null; default:0;   type:decimal(12,2)" json:"before_balance"`  //操作前余额
	AfterBalance  float64 `gorm:"not null; default:0;   type:decimal(12,2)" json:"after_balance"`   //操作后余额
	CreateTime    string  `gorm:"null;  type:datetime; json:"create_time"`                          //创建时间
	UniqueOpId    string  `gorm:"not null; default:'';  type:varchar(100); json:"unique_op_id"`     //操作唯一码
	RelateId      int64   `gorm:"not null; default:0;   type:int; json:"relate_id"`                 //支付操作ID
	BalanceType   string  `gorm:"not null; default:'';  type:varchar(100); json:"balance_type"`     //操作字段类型
	CreateAt      int64   `gorm:"not null; default:0;   type:int; json:"create_at"`                 //时间戳
	Remark1       string  `gorm:"not null; default:''; type:varchar(128); column:remark1" json:"-"` //备用1
	Remark2       string  `gorm:"not null; default:''; type:varchar(128); column:remark2" json:"-"` //备用2
}

func (Balance) TableName() string {
	return "cBalance"
}

func (BalanceLog) TableName() string {
	return "cBalanceLog"
}

const (
	SQL_InsertRedPacketParams = "INSERT INTO `cRedPacketParams` (`rp_id`,`red_packet_no`,`red_packet_money`,`status`,`lucky`) VALUES (%d,%d,%.2f,%d,%d)"
	SQL_QueryBalanceByUserId  = "SELECT * FROM cBalance  WHERE `user_id` = %d AND `status` = 1 LIMIT 1  "
	//SQL_InsertSpending       = "INSERT INTO `cSpending` (`wedding_id`,`user_id`,`operate_type`,`money`,`create_time`,`status`,`op_unique_no`,`create_at`,`remark1`) VALUES (%d,%d,%d,%.2f,'%s',1,'%s',%d,'%s')  "
)

//抢红包
func GrabRedPacket(user_id, rp_id, rp_params_id int64, money float64) error {
	tx := db.Begin()
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()
	now := time.Now()
	err := updateBalance(tx, user_id, money, 0, now, utils.GetGuid())
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = tx.Exec(fmt.Sprintf(GRAB_RED_PACKET, rp_id)).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = tx.Exec(fmt.Sprintf(GRAB_RED_PACKET_PARAMS, user_id, now.Unix(), rp_id, rp_params_id)).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}
	tx.Commit()
	tx = nil
	return nil
}

func QueryBalanceByUserId(userid int64) (*Balance, error) {
	balance := &Balance{}
	err := db.Raw(fmt.Sprintf(SQL_QueryBalanceByUserId, userid)).Scan(balance).Error
	if err != nil {
		log.Println(err.Error())
	}
	return balance, err
}

//user_id
func GenRedPacket(userId int64, useNum utils.RedFlash, req *config.Req_GenRedPacket, isBalanceType bool) (map[string]interface{}, string, error) {
	bill_no := ""
	tx := db.Begin()
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()
	resultMap := make(map[string]interface{})
	now := time.Now()             //当前时间
	opUniqueNo := utils.GetGuid() //当前操作的uniqueNo

	//spending 表 操作
	spending := Spending{WeddingId: req.Data.WeddingId,
		UserId:      userId,
		FounderId:   0,
		OperateType: 2,
		Money:       useNum.Money,
		CreateTime:  now.Format("2006-01-02 15:04:05"),
		Status:      1,
		OpUniqueNo:  opUniqueNo,
		CreateAt:    now.Unix(),
		Remark1:     req.Data.Wish}
	err := tx.Create(&spending).Error
	if err != nil {
		log.Println("insert into spending err：", err)
		return resultMap, bill_no, ERROR_DB_ACTION
	}

	var red_packet_status int64 = 0
	//如果是余额支持要修改这个
	if isBalanceType {
		err := updateBalance(tx, userId, -useNum.Money, spending.Id, now, opUniqueNo)
		if err != nil {
			return nil, bill_no, err
		}
		red_packet_status = 1
	} else {
		bill_no, err = InsertRecharge(userId, useNum.Money, opUniqueNo, now.Unix())
		if err != nil {
			return nil, bill_no, err
		}
	}

	//生成红包
	bQuestion, err := json.Marshal(req.Data.Question)
	if err != nil {
		log.Println(err.Error())
		return resultMap, bill_no, ERROR_DB_ACTION
	}
	redPacket := RedPacket{Guid: utils.GetGuid(),
		UserId:         userId,
		RoomId:         req.Data.RoomId,
		RedPacketNum:   req.Data.RedPacketNum,
		RedPacketMoney: useNum.Money,
		RedPacketType:  req.Data.RedPacketType,
		GetNum:         0,
		Status:         red_packet_status,
		EndStatus:      0,
		CreateAt:       now.Unix(),
		Question:       string(bQuestion),
		Remark1:        req.Data.Wish}
	err = tx.Create(&redPacket).Error
	if err != nil {
		log.Println("inser into redpacket err:", err)
		return resultMap, bill_no, ERROR_DB_ACTION
	}

	//生成红包明细
	for i := 0; i < int(req.Data.RedPacketNum); i++ {
		isLuck := 0
		if i == int(useNum.IndexMax) {
			isLuck = 1
		}
		err = tx.Exec(fmt.Sprintf(SQL_InsertRedPacketParams,
			redPacket.RpId, int(i+1), useNum.ResultRedPacketData[i], red_packet_status, isLuck)).Error
		if err != nil {
			log.Println("inser into redpacketparams err:", err)
			return resultMap, bill_no, ERROR_DB_ACTION
		}
	}
	//保存入redis数据库
	_, err = GetRedisDB().Do("SET", fmt.Sprintf("%s%d", config.REDIS_REDPACK_USER, redPacket.RpId), "")
	if err != nil {
		if strings.Index(err.Error(), "connection timed out") > 0 {
			ConnectRedis()
			_, err = GetRedisDB().Do("SET", fmt.Sprintf("%s%d", config.REDIS_REDPACK_USER, redPacket.RpId), "")
		}
		if err != nil {
			log.Println("redis set:", err)
			return resultMap, bill_no, ERROR_DB_ACTION
		}
	}
	//格式：1_3.59;2_7.93;3_3.48;
	redis_redpack := ""
	for k, v := range useNum.ResultRedPacketData {
		redis_redpack = fmt.Sprintf("%s%d_%.2f;", redis_redpack, k+1, v)
	}
	_, err = GetRedisDB().Do("SET", fmt.Sprintf("%s%d", config.REDIS_REDPACK, redPacket.RpId), redis_redpack)
	if err != nil {
		log.Println("redis set failed:", err)
		return resultMap, bill_no, ERROR_DB_ACTION
	}
	_bQuestion, err := json.Marshal(req.Data.Question)
	if err != nil {
		log.Println("redis set failed:", err)
		return resultMap, bill_no, ERROR_DB_ACTION
	}
	//返回Map
	resultMap = map[string]interface{}{"rp_id": redPacket.RpId,
		"guid":     redPacket.Guid,
		"wish":     req.Data.Wish,
		"question": string(_bQuestion)}
	tx.Commit()
	tx = nil
	return resultMap, bill_no, nil
}

//修改balance 后 insert into balancelog   balance 大于0 表示收入  小于0 表示消费
func updateBalance(tx *gorm.DB, userid int64, balance float64, spendingid int64, now time.Time, opUniqueNo string) error {
	//查询修改前 balance信息
	beforeBalance := &Balance{}
	err := tx.Raw(fmt.Sprintf(SQL_QueryBalanceByUserId, userid)).Scan(beforeBalance).Error
	if err != nil {
		log.Println(err)
		return ERROR_NOT_FIND_USER
	}

	//update balance
	update_field_name := "reward_num"
	if balance < 0 {
		//消费
		if beforeBalance.Balance < -balance {
			return errors.New("余额不足")
		}
		err = tx.Table("cBalance").Where(&Balance{UserId: userid, Status: 1}).
			Updates(map[string]interface{}{"balance": gorm.Expr("balance  - ?", -balance),
				"reward_num": gorm.Expr("reward_num  + ?", -balance),
				"rand_str":   utils.GetRandStr(),
				"sign":       gorm.Expr(SIGN)}).Error
		if err != nil {
			log.Println(err)
			return ERROR_DB_ACTION
		}
	} else {
		//收入
		err = tx.Table("cBalance").Where(&Balance{UserId: userid, Status: 1}).
			Updates(map[string]interface{}{"balance": gorm.Expr("balance  + ?", balance),
				"recharge_num": gorm.Expr("recharge_num  + ?", balance),
				"rand_str":     utils.GetRandStr(),
				"sign":         gorm.Expr(SIGN)}).Error
		if err != nil {
			log.Println(err)
			return ERROR_DB_ACTION
		}
		update_field_name = "recharge_num"
	}
	//查询修改后 balance信息
	afterBalance := &Balance{}
	err = tx.Raw(fmt.Sprintf(SQL_QueryBalanceByUserId, userid)).Scan(afterBalance).Error
	if err != nil {
		log.Println(err)
		return ERROR_NOT_FIND_USER
	}

	var before, after, _balance float64 = 0, 0, 0
	if update_field_name == "reward_num" {
		before, after, _balance = beforeBalance.RewardNum, afterBalance.RewardNum, -balance
	} else {
		before, after, _balance = beforeBalance.RechargeNum, afterBalance.RechargeNum, balance
	}
	//insert into balanceLog
	balanceLog_balance := BalanceLog{UserId: userid,
		BalanceId:     beforeBalance.Id,
		OperateType:   6,
		OperateValue:  balance,
		BeforeBalance: beforeBalance.Balance,
		AfterBalance:  afterBalance.Balance,
		CreateTime:    now.Format("2006-01-02 15:04:05"),
		UniqueOpId:    opUniqueNo,
		RelateId:      spendingid,
		BalanceType:   "balance",
		CreateAt:      now.Unix()}
	err = tx.Create(&balanceLog_balance).Error
	if err != nil {
		log.Println(err)
		return ERROR_DB_ACTION
	}

	balanceLog_reward_num := BalanceLog{UserId: userid,
		BalanceId:     beforeBalance.Id,
		OperateType:   6,
		OperateValue:  _balance,
		BeforeBalance: before,
		AfterBalance:  after,
		CreateTime:    now.Format("2006-01-02 15:04:05"),
		UniqueOpId:    opUniqueNo,
		RelateId:      spendingid,
		BalanceType:   update_field_name,
		CreateAt:      now.Unix()}
	err = tx.Create(&balanceLog_reward_num).Error
	if err != nil {
		log.Println(err)
		return ERROR_DB_ACTION
	}
	return nil
}
