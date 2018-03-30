package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
<<<<<<< HEAD
	"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/models"
	"github.com/caibinsong/wedding/utils"
=======
	"github.com/Amniversary/wedding-logic-redpacket/config"
	"github.com/Amniversary/wedding-logic-redpacket/models"
	"github.com/Amniversary/wedding-logic-redpacket/utils"
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	"gopkg.in/chanxuehong/wechat.v2/mch/core"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//微信创建红包
func WXGenRedPacket(w http.ResponseWriter, r *http.Request) {
	Response := &config.Response{Code: config.RESPONSE_ERROR}
	defer func() {
		EchoJson(w, http.StatusOK, Response)
	}()
	code := r.Header.Get("code")
	if code == "" {
		Response.Msg = "请先登录"
		return
	}
	//解析request中的数据
	req := &config.Req_GenRedPacket{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("request genRedPacket json decode err: %v", err)
		return
	}
	//还要判断答题红包问题不可为空
	err := req.CheckParameter()
	if err != nil {
		Response.Msg = err.Error()
		return
	}

	//红包金额数组
	redFlash, err := utils.GenRedPacket(req.Data.RedPacketType, float64(req.Data.RedPacketNum), req.Data.RedPacketMoney)
	if err != nil {
		Response.Msg = err.Error()
		return
	}
	log.Println(redFlash)
	//用户ID
	userid := GetUserId(r)
	if userid == 0 {
		Response.Msg = "用户ID异常"
		return
	}
	//把数据保存入库
	result, err := models.GenRedPacket(userid, redFlash, req, false)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "生成失败"
		return
	}
	////
	//生成roomsvr 广播信息
	roomMsg := map[string]interface{}{"rp_id": result["rp_id"], "type": req.Data.RedPacketType, "wish": result["wish"]}
	bRoomMsg, err := json.Marshal(roomMsg)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "生成广播失败"
		return
	}

	roomSvr := map[string]interface{}{"chatroomId": req.Data.RoomId,
		"weddingId": req.Data.WeddingId,
		"userId":    userid,
		"msgType":   4,
		"msg":       string(bRoomMsg)}
	bAttach, err := json.Marshal(roomSvr)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "生成广播失败"
		return
	}
	/////
	//
	rsp, err := NewWXRedPacket(result["rp_id"].(int64), result["guid"].(string), int64(redFlash.Money*100), code, string(bAttach))
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "生成失败"
		return
	}
	Response.Data = map[string]string{"appId": rsp["appid"],
		"nonceStr":  rsp["nonce_str"],
		"package":   "",
		"signType":  "MD5",
		"timeStamp": fmt.Sprint(time.Now().Unix()),
		"paySign":   rsp["sign"],
		"bill_no":   rsp["prepay_id"]}
	Response.Code = config.RESPONSE_OK
}

func NewWXRedPacket(rp_id int64, guid string, money int64, code, attach string) (map[string]string, error) {
<<<<<<< HEAD
	var req map[string]string = map[string]string{"appid": config.GetConfig().AppId,
		"attach":           fmt.Sprintf("%s%d;roomsvr=%s", config.ATTR_STR, rp_id, attach),
		"body":             "<![CDATA[微信支付充值]]>",
		"goods_tag":        "<![CDATA[微信支付充值]]></goods_tag>",
		"mch_id":           config.GetConfig().MchId,
		"nonce_str":        utils.GetMd5String(fmt.Sprintf("%d", time.Now().Unix())),
		"notify_url":       config.GetConfig().NotifyUrl,
		"openid":           code,
		"out_trade_no":     strings.Replace(guid, "-", "", -1),
		"spbill_create_ip": config.GetConfig().SpbillCreateIp,
		"total_fee":        fmt.Sprint(money),
		"trade_type":       "JSAPI"}

	client := core.NewClient(config.GetConfig().AppId, config.GetConfig().MchId, config.GetConfig().Key, nil)
=======
	var req map[string]string = map[string]string{"appid": config.APP_ID,
		"attach":           fmt.Sprintf("%s%d;roomsvr=%s", config.ATTR_STR, rp_id, attach),
		"body":             "<![CDATA[微信支付充值]]>",
		"goods_tag":        "<![CDATA[微信支付充值]]></goods_tag>",
		"mch_id":           config.MCH_ID,
		"nonce_str":        utils.GetMd5String(fmt.Sprintf("%d", time.Now().Unix())),
		"notify_url":       config.NOTIFY_URL,
		"openid":           code,
		"out_trade_no":     strings.Replace(guid, "-", "", -1),
		"spbill_create_ip": config.SPBILL_CREATE_IP,
		"total_fee":        fmt.Sprint(money),
		"trade_type":       "JSAPI"}

	client := core.NewClient(config.APP_ID, config.MCH_ID, config.KEY, nil)
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	response, err := client.PostXML("https://api.mch.weixin.qq.com/pay/unifiedorder", req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if response["result_code"] != "SUCCESS" {
		return nil, fmt.Errorf(response["return_msg"])
	}
	return response, nil
	/*
		map[
			result_code:SUCCESS
			prepay_id:wx2018032715431048c68a7d400519433148
			mch_id:1456793102
			nonce_str:CXZ97FTxhjuOZQ9Q
			appid:wx1e16505b46f55fc3
			sign:7F543FAEA02907787230C5158FF10257
			trade_type:JSAPI
			return_code:SUCCESS
			return_msg:OK
		] <nil>
	*/
}

func checkRequest(req *config.WXPayNotifyReq) bool {
<<<<<<< HEAD
	if req.Return_code == "SUCCESS" && req.Appid == config.GetConfig().AppId && req.Mch_id == config.GetConfig().MchId {
=======
	if req.Return_code == "SUCCESS" && req.Appid == config.APP_ID && req.Mch_id == config.MCH_ID {
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
		return true
	}
	log.Println("checkRequest err:", req.Return_code, req.Appid, req.Mch_id)
	return false
}

//pay_target=redpacket;rpid=1;roomsvr={fdsafa:aa}  解析出rpid 和 roomsvr广播的信息
func parsAttach(attach string) (int, string, error) {
	if !strings.HasPrefix(attach, config.ATTR_STR) {
		return 0, "", fmt.Errorf("Attach error : %s", attach)
	}
	attach = attach[len(config.ATTR_STR):]
	index := strings.Index(attach, ";")
	if index <= 0 {
		return 0, "", fmt.Errorf("Attach error : %s", attach)
	}
	str_rp_id := attach[:index]
	rp_id, err := strconv.Atoi(str_rp_id)
	if err != nil {
		log.Println(err.Error())
		return 0, "", err
	}
	index = strings.Index(attach, config.ROOM_SVR)
	if index <= 0 {
		return 0, "", fmt.Errorf("Attach error : %s", attach)
	}
	return rp_id, attach[index+len(config.ROOM_SVR):], nil
}

//
func CallBack(w http.ResponseWriter, r *http.Request) {
	req := &config.WXPayNotifyReq{}
	err := xml.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Println("解析HTTP Body格式到xml失败，原因!", err)
		EchoWXXML(w, http.StatusOK, "FAIL")
		return
	}

	if !checkRequest(req) {
		EchoWXXML(w, http.StatusOK, "FAIL")
		return
	}

	if !strings.HasPrefix(req.Attach, config.ATTR_STR) {
		log.Println("req.Attach:", req.Attach)
		EchoWXXML(w, http.StatusOK, "FAIL")
		return
	}
	rp_id, room_msg, err := parsAttach(req.Attach)
	if err != nil {
		log.Println(err.Error())
		EchoWXXML(w, http.StatusOK, "FAIL")
		return
	}
	if models.UpDateRedPacketStatus(int64(rp_id)) != nil {
		EchoWXXML(w, http.StatusOK, "FAIL")
		return
	}
	var roomMsg map[string]interface{} = make(map[string]interface{})
	err = json.Unmarshal([]byte(room_msg), roomMsg)
	if err != nil {
		log.Println(err.Error())
		EchoWXXML(w, http.StatusOK, "FAIL")
		return
	}
	log.Println(roomMsg)
	err = utils.NewHttpClient().RoomSvr(roomMsg)
	if err != nil {
		log.Println(err.Error())
		EchoWXXML(w, http.StatusOK, "FAIL")
		return
	}
	EchoWXXML(w, http.StatusOK, "SUCCESS")
}
