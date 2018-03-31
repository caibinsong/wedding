package controllers

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/models"
	"github.com/caibinsong/wedding/utils"
	"gopkg.in/chanxuehong/wechat.v2/mch/core"
	"log"
	"net/http"
	"regexp"
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
	result, err := models.GenRedPacket(userid, redFlash, genRedPacket, false)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "生成失败"
		return
	}
	////
	//生成roomsvr 广播信息
	// roomMsg := map[string]interface{}{"rp_id": result["rp_id"], "type": req.Data.RedPacketType, "wish": result["wish"]}
	// bRoomMsg, err := json.Marshal(roomMsg)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	Response.Msg = "生成广播失败"
	// 	return
	// }

	// roomSvr := map[string]interface{}{"chatroomId": req.Data.RoomId,
	// 	"weddingId": req.Data.WeddingId,
	// 	"userId":    userid,
	// 	"msgType":   4,
	// 	"msg":       string(bRoomMsg)}
	// bAttach, err := json.Marshal(roomSvr)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	Response.Msg = "生成广播失败"
	// 	return
	// }
	attach := ToSimpleAttach(result["rp_id"].(int64), req.Data.Params.RedPacketType, result["wish"].(string), req.Data.Params.RoomId, req.Data.Params.WeddingId, userid, 4)
	/////
	//
	rsp, err := NewWXRedPacket(result["rp_id"].(int64), result["guid"].(string), int64(redFlash.Money*100), userinfo.Data.OpenId, attach)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "生成失败"
		return
	}
	Response.Data = map[string]string{"appId": rsp["appid"],
		"nonceStr":  rsp["nonce_str"],
		"package":   fmt.Sprintf("prepay_id=%s", rsp["prepay_id"]),
		"signType":  "MD5",
		"timeStamp": fmt.Sprint(time.Now().Unix()),
		"paySign":   rsp["sign"],
		"bill_no":   "test_1234567"}
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

func NewWXRedPacket(rp_id int64, guid string, money int64, code, attach string) (map[string]string, error) {
	var req map[string]string = map[string]string{"appid": config.GetConfig().AppId,
		"attach":           attach,
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
	if req.Return_code == "SUCCESS" && req.Appid == config.GetConfig().AppId && req.Mch_id == config.GetConfig().MchId {
		return true
	}
	log.Println("checkRequest err:", req.Return_code, req.Appid, req.Mch_id)
	return false
}

// //pay_target=redpacket;rpid=1;roomsvr={fdsafa:aa}  解析出rpid 和 roomsvr广播的信息
// func parsAttach(attach string) (int, string, error) {
// 	if !strings.HasPrefix(attach, config.ATTR_STR) {
// 		return 0, "", fmt.Errorf("Attach error : %s", attach)
// 	}
// 	attach = attach[len(config.ATTR_STR):]
// 	index := strings.Index(attach, ";")
// 	if index <= 0 {
// 		return 0, "", fmt.Errorf("Attach error : %s", attach)
// 	}
// 	str_rp_id := attach[:index]
// 	rp_id, err := strconv.Atoi(str_rp_id)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return 0, "", err
// 	}
// 	index = strings.Index(attach, config.ROOM_SVR)
// 	if index <= 0 {
// 		return 0, "", fmt.Errorf("Attach error : %s", attach)
// 	}
// 	return rp_id, attach[index+len(config.ROOM_SVR):], nil
// }

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

	rp_id, room_msg, err := ToJsonAttach(req.Attach)
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
