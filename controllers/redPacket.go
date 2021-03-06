package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/models"
	"github.com/caibinsong/wedding/utils"
	"log"
	"net/http"
	"strconv"
)

//创建红包
func GenRedPacket(w http.ResponseWriter, r *http.Request) {
	Response := &config.Response{Code: config.RESPONSE_ERROR}
	defer func() {
		EchoJson(w, http.StatusOK, Response)
	}()
	//解析request中的数据
	req := &config.Req_GenRedPacket{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("request genRedPacket json decode err: %v", err)
		Response.Msg = "请求参数有误"
		return
	}
	//传入值合法性校验
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
	//用户ID
	userid := GetUserId(r)
	if userid == 0 {
		Response.Msg = "用户ID异常"
		return
	}
	//把数据保存入库
	result, _, err := models.GenRedPacket(userid, redFlash, req, true)
	if err != nil {
		Response.Msg = err.Error()
		return
	}
	_bQuestion, err := json.Marshal(req.Data.Question)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "广播失败"
		return
	}
	//发送roomsvr 广播
	roomMsg := map[string]interface{}{"rp_id": result["rp_id"], "type": req.Data.RedPacketType, "wish": result["wish"], "question": string(_bQuestion)}
	bRoomMsg, err := json.Marshal(roomMsg)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = "广播失败"
		return
	}

	roomSvr := map[string]interface{}{"chatroomId": req.Data.RoomId,
		"weddingId": req.Data.WeddingId,
		"userId":    userid,
		"msgType":   4,
		"msg":       string(bRoomMsg)}
	//	//config.RoomSvr_ServerName   AccessCtrl Broadcast
	err = utils.NewHttpClient().RoomSvr(config.RoomSvr_ServerName, config.RoomSvr_MethodName, roomSvr)
	if err != nil {
		Response.Msg = err.Error()
		return
	}
	Response.Data = result
	Response.Code = config.RESPONSE_OK
}

//抢红包
func GrabRedPacket(w http.ResponseWriter, r *http.Request) {
	Response := &config.Response{Code: config.RESPONSE_ERROR}
	defer func() {
		EchoJson(w, http.StatusOK, Response)
	}()
	//解析request中的数据
	req := &config.Req_RedPacket{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("request Req_RedPacket json decode err: %v", err)
		return
	}
	//用户ID
	userid := GetUserId(r)
	if userid == 0 {
		Response.Msg = "用户ID异常"
		return
	}

	err := models.GetRedPack(userid, req.Data.RpId)
	if err != nil {
		log.Println(err.Error())
		Response.Msg = err.Error()
	}
	// redPacket, err := models.FindRedPacketByRpId(req.Data.RpId)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	Response.Msg = "红包已经抢完！"
	// }
	Response.Data = map[string]interface{}{"rp_id": req.Data.RpId /*,"red_type": redPacket.RedPacketType*/}
	// speeding, err := models.QuerySpending(redPacket.UserId, redPacket.CreateAt)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	Response.Msg = "广播失败"
	// }
	// /////广播
	// data := map[string]interface{}{"rp_id": req.Data.RpId,
	// 	"red_type": redPacket.RedPacketType}
	// var msg map[string]interface{}
	// if Response.Msg == "" {
	// 	msg = map[string]interface{}{"code": 0}
	// } else {
	// 	msg = map[string]interface{}{"code": 1, "msg": Response.Msg}
	// }
	// roomSvr := map[string]interface{}{"chatroomId": redPacket.RoomId,
	// 	"weddingId": speeding.WeddingId,
	// 	"userId":    userid,
	// 	"data":      data,
	// 	"msg":       msg}

	// bRoomSvr, err := json.Marshal(map[string]interface{}{"msg": roomSvr})
	// if err != nil {
	// 	log.Println(err.Error())
	// 	Response.Msg = "广播失败"
	// 	return
	// }
	// content := map[string]interface{}{
	// 	"type":    2,
	// 	"content": string(bRoomSvr),
	// }
	// bContent, err := json.Marshal(content)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	Response.Msg = "广播失败"
	// 	return
	// }
	// body := map[string]interface{}{
	// 	"type":    "HLBUser",
	// 	"idList":  []int64{userid},
	// 	"content": string(bContent),
	// }
	models.AddAccessCtrWork(models.AccessCtr{RpId: req.Data.RpId, UserId: userid})
	Response.Code = config.RESPONSE_OK
}

//红包列表
func GetRedPacketInfo(w http.ResponseWriter, r *http.Request) {
	Response := &config.Response{Code: config.RESPONSE_ERROR}

	defer func() {
		EchoJson(w, http.StatusOK, Response)
	}()
	//解析request中的数据
	req := &config.Req_RedPacket{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("request Req_RedPacket json decode err: %v", err)
		return
	}
	res_data := &config.Res_GetRedPacketInfo{RpId: req.Data.RpId}
	//查询红包信息
	redPacket, err := models.FindRedPacketByRpId(req.Data.RpId)
	if err != nil {
		Response.Msg = "查询红包异常"
		return
	}

	//获取发红包人个人微信信息
	response, err := utils.NewHttpClient().GetWXUserInfoResponse(redPacket.UserId)
	if err != nil {
		Response.Msg = "个人信息获取失败"
		return
	}

	res_data.Red_Info = config.RedInfo{UserId: int(redPacket.UserId),
		NickName:       response.Data.NickName,
		Pic:            response.Data.Pic,
		Wish:           redPacket.Remark1,
		RedPacketNum:   redPacket.RedPacketNum,
		RedPacketMoney: redPacket.RedPacketMoney,
		RedPacketType:  redPacket.RedPacketType,
		Status:         redPacket.Status}

	//查询红包明细信息
	list, err := models.FindRedPacketParamsByRpId(req.Data.RpId)
	if err != nil {
		Response.Msg = "查询红包异常"
		return
	}
	userlist := make([]int, 0)
	for _, one_user := range list {
		userlist = append(userlist, int(one_user.UserId))
	}
	otherInfo := config.OtherInfo{Count: len(userlist)}
	otherInfo.List = make([]config.OtherList, 0)
	//获取明细微信信息列表
	if len(userlist) > 0 {
		response_list, err := utils.NewHttpClient().GetWXUserListResponse(userlist)
		if err != nil {
			Response.Msg = "个人信息获取失败"
			return
		}
		for _, one_userId := range list {
			for _, one := range response_list.Data {
				if fmt.Sprint(one_userId.UserId) == one.Id {
					add := config.OtherList{UserId: one_userId.UserId,
						NickName: one.NickName,
						Pic:      one.Pic,
						UpdateAt: one_userId.UpdateAt,
						Money:    one_userId.RedPacketMoney,
						Lucky:    one_userId.Lucky}
					otherInfo.List = append(otherInfo.List, add)
					break
				}
			}
		}
	}
	res_data.Other_Info = otherInfo
	Response.Data = res_data
	Response.Code = config.RESPONSE_OK
}

//是否可以抢这个红包
func CheckUserRedPacket(w http.ResponseWriter, r *http.Request) {
	Response := &config.Response{Code: config.RESPONSE_ERROR}
	defer func() {
		EchoJson(w, http.StatusOK, Response)
	}()
	//解析request中的数据
	req := &config.Req_RedPacket{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("request Req_RedPacket json decode err: %v", err)
		return
	}
	//用户ID
	userid := GetUserId(r)
	if userid == 0 {
		Response.Msg = "用户ID异常"
		return
	}

	status, err := models.CheckUserRedPacket(userid, req)
	if err != nil {
		Response.Msg = err.Error()
		return
	}
	Response.Data = status
	Response.Code = config.RESPONSE_OK
}

//获取header中 userid
func GetUserId(r *http.Request) int64 {
	sUserid := r.Header.Get("userid")
	userid, err := strconv.ParseInt(sUserid, 10, 64)
	if err != nil {
		return 0
	}
	return userid
}
