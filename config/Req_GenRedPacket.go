package config

import (
	"errors"
)

//请求上来的request对象
const (
	RED_PACKET_MAX       = 20000 //TODO: 单位 (分) 红包大小限制
	RED_PACKET_NUM_MAX   = 200
	ACTION_NAME          = "gen_red_packet"
	LUCKY_RED_PACKET     = 1 //手气红包
	GENERAL_RED_PACKET   = 2 //普通红包
	QUESTIONS_RED_PACKET = 3 //答题红包
)

//创建红包
type Req_GenRedPacket struct {
	ActionName string `json:"action_name"`
	Data       struct {
		WeddingId      int64   `json:"wedding_id"`
		RoomId         int64   `json:"room_id"`
		RedPacketMoney float64 `json:"red_packet_money"`
		RedPacketNum   int64   `json:"red_packet_num"`
		RedPacketType  int64   `json:"red_packet_type"`
		Wish           string  `json:"wish"`
		Question       struct {
			Title    string   `json:"title"`
			Key      int64    `json:"key"`
			Question []string `json:"question"`
		} `json:"question"`
	} `json:"data"`
}

//校验数据合法性（包括余额是否足够）
func (this *Req_GenRedPacket) CheckParameter() error {
	if this.ActionName != ACTION_NAME {
		return errors.New("ActionName参数不正确")
	}
	if this.Data.WeddingId <= 0 || this.Data.RoomId <= 0 || this.Data.RedPacketMoney <= 0 || this.Data.RedPacketNum <= 0 {
		return errors.New("参数不正确")
	}
	//红包个数
	if this.Data.RedPacketNum <= 0 || this.Data.RedPacketNum > RED_PACKET_NUM_MAX {
		return errors.New("红包个数不可超过200")
	}
	//红包类型与红包大小判断
	if this.Data.RedPacketType == GENERAL_RED_PACKET {
		//普通红包
		if this.Data.RedPacketMoney <= 0 || this.Data.RedPacketMoney > 200 {
			return errors.New("单个红包不能大于200元")
		}
	} else if this.Data.RedPacketType == LUCKY_RED_PACKET || this.Data.RedPacketType == QUESTIONS_RED_PACKET {
		//手气红包和答题红包
		//单个红包不可小于0.01
		if this.Data.RedPacketMoney < (float64(this.Data.RedPacketNum) * 0.01) {
			return errors.New("单个红包不能小于0.01元")
		}
		//单个红包不可大于200
		if this.Data.RedPacketMoney*100 > (float64(this.Data.RedPacketNum) * RED_PACKET_MAX) {
			return errors.New("单个红包不能大于200元")
		}
	} else {
		return errors.New("红包类型不正确")
	}

	//答题红包
	if this.Data.RedPacketType == QUESTIONS_RED_PACKET {
		if this.Data.Question.Title == "" {
			return errors.New("答题红包，题目不可为空")
		}
		if len(this.Data.Question.Question) < 2 {
			return errors.New("答题红包，答案不可少于2个")
		}
		if this.Data.Question.Key < 0 || this.Data.Question.Key >= int64(len(this.Data.Question.Question)) {
			return errors.New("正确答案选择有误")
		}
	}
	return nil
}

//抢红包接口 获取红包详情列表 判断是否可以抢红包
type Req_RedPacket struct {
	ActionName string `json:"action_name"`
	Data       struct {
		RpId int64 `json:"rp_id"`
	}
}
