package controllers

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/models"
	"github.com/caibinsong/wedding/utils"
	"github.com/chanxuehong/rand"
	"gopkg.in/chanxuehong/wechat.v2/mch/core"
	"gopkg.in/chanxuehong/wechat.v2/mch/pay"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//微信创建红包
func WXGenRedPacket(w http.ResponseWriter, r *http.Request) {
	Response := &config.Response{Code: config.RESPONSE_ERROR}
	defer func() {
		EchoJson(w, http.StatusOK, Response)
	}()
	//用户ID
	userid := GetUserId(r)
	if userid == 0 {
		Response.Msg = "用户ID异常"
		return
	}
	userinfo, err := utils.NewHttpClient().GetWXUserInfoResponse(userid)
	if err != nil {
		Response.Msg = err.Error()
		return
	}
	log.Println(userinfo)
	//解析request中的数据
	req := &config.Req_WXGenRedPacket{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("request genRedPacket json decode err: %v", err)
		return
	}
	//还要判断答题红包问题不可为空
	err = req.CheckParameter()
	if err != nil {
		Response.Msg = err.Error()
		return
	}

	//红包金额数组
	redFlash, err := utils.GenRedPacket(req.Data.Params.RedPacketType, float64(req.Data.Params.RedPacketNum), req.Data.Params.RedPacketMoney)
	if err != nil {
		Response.Msg = err.Error()
		return
	}
	log.Println(redFlash)
	var genRedPacket *config.Req_GenRedPacket = &config.Req_GenRedPacket{Data: req.Data.Params}
	log.Println(genRedPacket)
	//把数据保存入库
	result, bill_no, err := models.GenRedPacket(userid, redFlash, genRedPacket, false)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "生成失败"
		return
	}

	attach := ToSimpleAttach(result["rp_id"].(int64), req.Data.Params.RedPacketType, result["wish"].(string), req.Data.Params.RoomId, req.Data.Params.WeddingId, userid, 4)
	//
	rsp, paySign, nonceStr, timestamp, err := NewWXRedPacket(result["rp_id"].(int64), result["guid"].(string), int64(redFlash.Money*100),
		userinfo.Data.OpenId, attach)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "生成失败"
		return
	}
	log.Println(bill_no, rsp)
	Response.Data = map[string]string{"appId": config.GetConfig().AppId,
		"signType":  "MD5",
		"total_fee": fmt.Sprint(int64(redFlash.Money * 100)),
		"bill_no":   bill_no,
		"package":   fmt.Sprintf("prepay_id=%s", rsp.PrepayId),
		"timeStamp": timestamp,
		"nonceStr":  nonceStr,
		"paySign":   paySign,
	}
	Response.Code = config.RESPONSE_OK
}

func ToSimpleAttach(rpid, rptype int64, wish string, chatroomId, weddingId, userId, msgType int64) string {
	return fmt.Sprintf("ri=%d;rt=%d;wh=%s;ci=%d;wi=%d;ui=%d;mt=%d;",
		rpid, rptype, wish, chatroomId, weddingId, userId, msgType)
}

func ToJsonAttach(msg string) (int, string, error) {
	reg := regexp.MustCompile("ri=(.*?);rt=(.*?);wh=(.*?);ci=(.*?);wi=(.*?);ui=(.*?);mt=(.*?);")
	arr := reg.FindStringSubmatch(msg)
	if len(arr) != 8 {
		return 0, "", errors.New("数据不正确")
	}
	rp_id, err := strconv.Atoi(arr[1])
	if err != nil {
		return 0, "", errors.New("数据不正确")
	}
	rst := fmt.Sprintf(`{"chatroomId":%s,"msgType":%s,"userId":%s,"weddingId":%s,"msg":"{\"rp_id\":%s,\"type\":%s,\"wish\":\"%s\"}"}`,
		arr[4], arr[7], arr[6], arr[5], arr[1], arr[2], arr[3])
	return rp_id, rst, nil
}

func NewWXRedPacket(rp_id int64, guid string, money int64, openid, attach string) ( /*map[string]string*/ *pay.UnifiedOrderResponse, string, string, string, error) {
	nonceStr := string(rand.NewHex())
	uor := &pay.UnifiedOrderRequest{
		DeviceInfo:     "WEB",
		Detail:         "<![CDATA[微信支付充值]]>",
		Attach:         attach,
		Body:           "<![CDATA[微信支付充值]]>",
		GoodsTag:       "<![CDATA[微信支付充值]]></goods_tag>",
		NonceStr:       nonceStr,
		NotifyURL:      config.GetConfig().NotifyUrl,
		OpenId:         openid,
		OutTradeNo:     strings.Replace(guid, "-", "", -1),
		SpbillCreateIP: config.GetConfig().SpbillCreateIp,
		TotalFee:       money,
		TradeType:      "JSAPI",
	}
	client := core.NewClient(config.GetConfig().AppId, config.GetConfig().MchId, config.GetConfig().Key, nil)
	rsp, err := pay.UnifiedOrder2(client, uor)
	if err != nil {
		log.Println("mch pay unified order error: ", err)
		return nil, "", "", "", err
	}

	timestamp := fmt.Sprintf("prepay_id=%s", rsp.PrepayId)
	PaySign := core.JsapiSign(config.GetConfig().AppId, timestamp, nonceStr,
		timestamp,
		"MD5",
		config.GetConfig().Key)
	log.Println(rsp, err)
	return rsp, PaySign, nonceStr, timestamp, err
}

func checkRequest(req *config.WXPayNotifyReq) bool {
	if req.Return_code == "SUCCESS" && req.Appid == config.GetConfig().AppId && req.Mch_id == config.GetConfig().MchId {
		return true
	}
	log.Println("checkRequest err:", req.Return_code, req.Appid, req.Mch_id)
	return false
}

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
	log.Println(req.Out_trade_no)
	log.Println(req)
	rp_id, room_msg, err := ToJsonAttach(req.Attach)
	if err != nil {
		log.Println(err.Error())
		EchoWXXML(w, http.StatusOK, "FAIL")
		return
	}
	if models.UpDateRedPacketStatus(int64(rp_id), req.Transaction_id) != nil {
		EchoWXXML(w, http.StatusOK, "FAIL")
		return
	}
	log.Println(room_msg)
	var roomMsg map[string]interface{} = make(map[string]interface{})
	log.Println(room_msg)
	err = json.Unmarshal([]byte(room_msg), &roomMsg)
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
