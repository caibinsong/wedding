package main

import (
	//"crypto/rand"
	//"github.com/Amniversary/wedding-logic-redpacket/business"
	//"github.com/caibinsong/wedding/config"
	//"github.com/Amniversary/wedding-logic-redpacket/models"
	//"encoding/json"
	//"github.com/Amniversary/wedding-logic-redpacket/utils"
	"log"
	//"math/big"
	//"github.com/jinzhu/gorm"
	//"time"
	//"errors"
	//"fmt"
	"github.com/caibinsong/wedding/controllers"
	//"regexp"
	//"strconv"
)

// roomMsg := map[string]interface{}{"rp_id": result["rp_id"], "type": req.Data.RedPacketType, "wish": result["wish"]}
// 	bRoomMsg, err := json.Marshal(roomMsg)
// 	if err != nil {
// 		log.Println(err.Error())
// 		Response.Msg = "生成广播失败"
// 		return
// 	}
// 	roomSvr := map[string]interface{}{"chatroomId": req.Data.RoomId,
// 		"weddingId": req.Data.WeddingId,
// 		"userId":    userid,
// 		"msgType":   4,
// 		"msg":       string(bRoomMsg)}

// "chatroomId": req.Data.RoomId,
// "weddingId": req.Data.WeddingId,
// "userId":    userid,
// "msgType":   4,
// "msg":       string(bRoomMsg)

//pay_target=redpacket;rpid=7;
//roomsvr=
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	attach := controllers.ToSimpleAttach("ri=(.*?);rt=(.*?);wh=(.*?);ci=(.*?);wi=1;ui=1;mt=1;")
	log.Println(attach, len(attach))
	log.Println(ToJsonAttach(attach))
	//config.InitConfig()
	//log.Println(config.GetConfig())
	// log.Println("start")
	// models.InitDataBase()
	// //log.Println(models.GetRedPacketInfo(18))
	// b0, err := models.QueryBalanceByUserId(1)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// log.Println(b0)
	// b1, err := models.QueryBalanceByUserId(2)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// log.Println(b1)
	/*
		Id         int64  `gorm:"primary_key" json:"id"`
		UserId     int64  `gorm:"not null; default:0;  type:int" json:"user_id"`
		RealName   int64  `gorm:"not null; default:0;  type:int" index" json:"real_name"`
		Phone      int64  `gorm:"not null; default:0;  type:int" index" json:"phone"`
		CreateTime string `gorm:"null;  type:datetime;" json:"create_time"`
		Remark1    string `gorm:"not null; default:''; type:varchar(128); column:remark1" json:"-"`
		Remark2
	*/
	//user := models.Client{UserId: 1, RealName: "蔡斌松111", Phone: "15267093345",
	//	CreateTime: time.Now().Format("2006-01-02 15:04:05")}

	//insert
	//log.Println(models.GetDBObject().Create(&user))
	//log.Println(user.Id)

	//find := models.Client{}
	//select
	//log.Println(models.GetDBObject().Where(&models.Client{Id: 1111}).First(&find))
	//log.Println(find, find.UserId)
	//update
	// c := models.Client{}
	// models.GetDBObject().Model(&c).Where(&models.Client{Id: 1, RealName: "蔡斌松"}).Updates(models.Client{RealName: "hello", Phone: "111111"})
	// log.Println(c.Id)
	//models.GetDBObject().Table("cBalance").Where(&models.Balance{UserId: 1}).Updates(map[string]interface{}{"balance": gorm.Expr("balance  - ?", 100), "reward_num": gorm.Expr("reward_num  + ?", 100)})
	//
	// models.InitDataBase()
	// log.Println(models.GetRedPack(1, 11))
	// var header map[string]string = map[string]string{"ServerName": "wedding", "MethodName": "getUserInfo", "userId": "844"}
	// var request *config.GetWXUserInfo = &config.GetWXUserInfo{ActionName: "get_user_info", Data: "Data"}
	// var response *config.WXUserInfoResponse = &config.WXUserInfoResponse{}

	// body, err := utils.NewHttpClient().Post("https://access.hunlibaoapp.com/socket/response.do", header, request)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// err = json.Unmarshal(body, &response)
	// if err != nil {
	// 	log.Println("json decode error: %v [%s]", err, string(body))
	// 	return
	// }
	// if response.Code != 0 {
	// 	log.Println("login robot result code error: %d %s", response.Code, response.Msg)
	// 	return
	// }
	// log.Println(response)
}
