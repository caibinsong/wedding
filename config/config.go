package config

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	ONE_DAY            = 24 * 60 * 60 //一天的秒数
	WX_GetUserInfo     = "getUserInfo"
	WX_GetUserList     = "getUserList"
	RoomSvr_ServerName = "Chatroom"
	RoomSvr_MethodName = "SendMessage"
	RoomSvr_Broadcast  = "Broadcast"
	RESPONSE_OK        = 0
	RESPONSE_ERROR     = 1
	ERROR_MSG          = "系统错误"
	//redis Key前缀
	REDIS_REDPACK      = "redpack_"      //+ rp_id
	REDIS_REDPACK_USER = "redpack_user_" //+ rp_id
	WX_RETURN          = `<xml>
						  <return_code><![CDATA[%s]]></return_code>
						  <return_msg><![CDATA[%s]]></return_msg>
						</xml>`
	ATTR_STR = "pay_target=redpacket;rpid="
	ROOM_SVR = "roomsvr="
)

type MySql struct {
	DBName string `xml:"dbname"`
	Host   string `xml:"host"`
	User   string `xml:"user"`
	Pass   string `xml:"pass"`
	Debug  string `xml:"debug"`
}

type Config struct {
	ServerName     string `xml:"servername"`
	Mysql          MySql  `xml:"mysql"`
	Redis          string `xml:"redis"`
	RedisPassWord  string `xml:"redispassword"`
	WXUserInfoUrl  string `xml:"wxuserinfourl"`
	WXUserListUrl  string `xml:"wxuserlisturl"`
	RoomSvrUrl     string `xml:"roomsvrurl"`
	AppId          string `xml:"appid"`
	MchId          string `xml:"mchid"`
	Key            string `xml:"key"`
	SpbillCreateIp string `xml:"spbillcreateip"`
	NotifyUrl      string `xml:"notifyurl"`
}

var config *Config = nil

func InitConfig() {

	config = &Config{}
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Panic(err.Error())
	}
	config_path := filepath.Join(filepath.Dir(path), "config.xml")
	log.Println(config_path)
	byts, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Panic("读取config.xml出错", err.Error())
	}
	err = xml.Unmarshal(byts, &config)
	if err != nil {
		log.Panic("解析config.xml出错", err.Error())
	}
}

func GetConfig() *Config {
	return config
}

// package config

// const (
// 	ServerName = "wedding" //服务名

// 	ONE_DAY = 24 * 60 * 60

// 	WX_GetUserInfo     = "getUserInfo"
// 	WX_GetUserList     = "getUserList"
// 	RoomSvr_ServerName = "Chatroom"
// 	RoomSvr_MethodName = "SendMessage"

// 	RESPONSE_OK      = 0
// 	RESPONSE_ERROR   = 1
// 	ERROR_MSG        = "系统错误"
// 	WX_USER_INFO_URL = "https://access.hunlibaoapp.com/socket/response.do"
// 	WX_USER_LIST_URL = "https://access.hunlibaoapp.com/socket/response.do"
// 	RoomSvr_URL      = "https://access.hunlibaoapp.com/socket/response.do"

// 	//redis Key前缀
// 	REDIS_REDPACK      = "redpack_"      //+ rp_id
// 	REDIS_REDPACK_USER = "redpack_user_" //+ rp_id

// 	//微信支付
// 	APP_ID           = "wx1e16505b46f55fc3"
// 	MCH_ID           = "1456793102"
// 	NOTIFY_URL       = "https://logic.hunlibaoapp.com/wxpay/wxpay_notify"
// 	SPBILL_CREATE_IP = "172.17.0.13"
// 	KEY              = "Ej1OlMdS54J8AfsZcmZGfYufNqYPSXuH"
// 	WX_RETURN        = `<xml>
// 						  <return_code><![CDATA[%s]]></return_code>
// 						  <return_msg><![CDATA[%s]]></return_msg>
// 						</xml>`

// 	ATTR_STR = "pay_target=redpacket;rpid="
// 	ROOM_SVR = "roomsvr="

// 	//my数据库
// 	DBName = "cAuth"     //wedding_card
// 	USER   = "root"      //root
// 	PASS   = "123456"    //tkC42cwy2U3SQwHw
// 	HOST   = "127.0.0.1" //172.17.0.5
// 	DEBUG  = "dev"       //prod

// 	//redis数据库连接
// 	REDIS_CONNECT = "127.0.0.1:6379"
// )

// /*
// 	//BalanceLog - OperateType
// 	1 //收到礼金 加余额
// 	2 //收到礼金 增加累计收入
// 	3 //提现 扣余额
// 	4 //发出礼金 增加累计支出
// 	5 //提现 增加累计提现
// 	6 //发红包 扣余额
// 	7 //发红包 加累计支出
// 	8 //抢红包 加余额
// 	9 //抢红包 加累计收入
// 	10 //退红包 增加余额
// 	11 //收礼物 增加余额
// 	12 //收礼物 增加累计收入
// 	13 //送礼物 减少余额
// 	14 //送礼物 增加累计支出
// 	15 //发送礼金 余额减少
// 	16 //提现 退提现金额 加余额
// 	//reward - OperateType
// 	1//.送礼金
// 	2//.发红包
// 	3//.送礼物
// */
