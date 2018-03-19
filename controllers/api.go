package controllers

import (
	"net/http"
	"github.com/Amniversary/wedding-logic-redpacket/config"
	"fmt"
	"encoding/json"
	"time"
	"github.com/Amniversary/wedding-logic-redpacket/models"
	"log"
)

func init() {
	models.InitDataBase()
	http.HandleFunc("/rpc", RunRpc)
}

func Run() {
	http.ListenAndServe(":5501", nil)
}

func RunRpc(w http.ResponseWriter, r *http.Request) {
	res := &config.Response{Code: config.RESPONSE_OK}
	if r.Method != "POST" {
		log.Printf("Method not be Post Request [%s]\n", r.Method)
		EchoJson(w, http.StatusOK, res)
		return
	}
	serverName := r.Header.Get("ServerName")
	if serverName != config.ServerName {
		log.Printf("ServerName: [%s]  request -> ServerName: [%s] Method: [%s]\n", config.ServerName, serverName, r.Method)
		EchoJson(w, http.StatusOK, res)
		return
	}

	methodName := r.Header.Get("MethodName")
	start := time.Now()
	defer func() {
		log.Printf("Request MethodName: [%s], Rtime[%v]\n", methodName, time.Now().Sub(start))
	}()
	switch methodName {
	case config.GEN_RED_PACKET:
		GenRedPacket(w, r)

	default:
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
