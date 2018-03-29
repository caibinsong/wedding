package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Amniversary/wedding-logic-redpacket/config"
	"github.com/Amniversary/wedding-logic-redpacket/models"
	"io"
	"log"
	"net/http"
	"time"
)

type MethodFunc func(w http.ResponseWriter, r *http.Request)

//方法映射Map
var MethodMap map[string]MethodFunc = map[string]MethodFunc{"genRedPacket": GenRedPacket, //创建红包
	"grabRedPacket":      GrabRedPacket,      //抢红包
	"getRedPacketInfo":   GetRedPacketInfo,   //红包列表
	"checkUserRedPacket": CheckUserRedPacket, //是否可以抢这个红包
	"WXGenRedPacket":     WXGenRedPacket}     //微信生成红包

func init() {
	//数据库初始化
	models.InitDataBase()
	http.HandleFunc("/rpc", RunRpc)
	http.HandleFunc("/wechat_callback", CallBack)
}

func Run() {
	http.ListenAndServe(":5501", nil)
}

//请求入口
func RunRpc(w http.ResponseWriter, r *http.Request) {
	res := &config.Response{Code: config.RESPONSE_OK}

	//请求必须为POST
	if r.Method != "POST" {
		log.Printf("Method not be Post Request [%s]\n", r.Method)
		EchoJson(w, http.StatusOK, res)
		return
	}

	//服务名必须一致
	serverName := r.Header.Get("ServerName")
	if serverName != config.ServerName {
		log.Printf("ServerName: [%s]  request -> ServerName: [%s] Method: [%s]\n", config.ServerName, serverName, r.Method)
		EchoJson(w, http.StatusOK, res)
		return
	}

	//userid，如果没有必须先授权
	userid := r.Header.Get("userid")
	if userid == "" {
		res.Code = 10002
		res.Msg = "用户未授权, 请授权小程序"
		EchoJson(w, http.StatusOK, res)
		return
	}
	methodName := r.Header.Get("MethodName")

	start := time.Now()
	defer func() {
		log.Printf("Request MethodName: [%s], Rtime[%v]\n", methodName, time.Now().Sub(start))
	}()
	methodFunc, ok := MethodMap[methodName]
	if ok {
		methodFunc(w, r)
	} else {
		res.Code = 1
		res.Msg = fmt.Sprintf("Can't find the interface: [%s]", methodName)
		EchoJson(w, http.StatusOK, res)
	}
}

// TODO @ 输出Json数据
func EchoJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type,servername,methodname,userid,msgid")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// TODO @ 输出Json数据
func EchoWXXML(w http.ResponseWriter, status int, return_code string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type,servername,methodname,userid,msgid")
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	io.WriteString(w, fmt.Sprintf(config.WX_RETURN, return_code))
}
