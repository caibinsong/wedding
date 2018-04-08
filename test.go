package main

import (
	//"crypto/rand"
	//"github.com/Amniversary/wedding-logic-redpacket/business"
	//"github.com/caibinsong/wedding/config"
	//"github.com/caibinsong/wedding/config"
	//"github.com/caibinsong/wedding/models"
	//"encoding/json"
	"fmt"
	"github.com/caibinsong/wedding/utils"
	"log"
	"time"
	//"math/big"
	//"github.com/jinzhu/gorm"
	//"time"
	//"errors"
	//"fmt"
	//"github.com/caibinsong/wedding/controllers"
	//"regexp"
	//"strconv"
	"io/ioutil"
	"strings"
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
func post(id int) {
	var hear = map[string]string{"ServerName": "weddingRedPacket",
		"MethodName":   "grabRedPacket",
		"userid":       fmt.Sprint(id),
		"Content-Type": "application/json"}
	var data = map[string]interface{}{"rp_id": 198}
	var request = map[string]interface{}{"action_name": "grab_red_packet",
		"data": data}
	r, err := utils.NewHttpClient().Post("http://182.254.247.115:5501/rpc", hear, request)
	//r, err := utils.NewHttpClient().Post("http://127.0.0.1:5501/rpc", hear, request)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()
	a, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if strings.Index(string(a), `"code":0`) > 0 {
		log.Println(id)
	}

}
func main() {
	for i := 100; i < 750; i++ { //843
		go post(i)
	}
	time.Sleep(time.Second * 200)
	//log.SetFlags(log.Lshortfile | log.LstdFlags)
	// go models.StartAccessCtrWork("aaaa", "bbbb")
	// log.Println(1)
	// for i := 0; i <= 10; i++ {
	// 	a := map[string]interface{}{"a": i}
	// 	models.AddAccessCtrWork(a)
	// }
	// log.Println(2)
	// time.Sleep(time.Second * 10)
	// attach := controllers.ToSimpleAttach("ri=(.*?);rt=(.*?);wh=(.*?);ci=(.*?);wi=1;ui=1;mt=1;")
	// log.Println(attach, len(attach))
	// log.Println(ToJsonAttach(attach))
	//config.InitConfig()
	//log.Println(config.GetConfig())
	// log.Println("start")
	//num := 10
	// log.Println("start")
	// config.InitConfig()
	// models.InitDataBase()
	// start := time.Now()
	// log.Println(models.QueryBalanceByUserId(1))
	// log.Println(time.Now().Sub(start))
	// start = time.Now()
	// log.Println(models.QueryBalanceByUserId1(1))
	// log.Println(time.Now().Sub(start))

	// start = time.Now()
	// log.Println(models.QuerySpending(1, 1522469859))
	// log.Println("insert ", time.Now().Sub(start))
	// start := time.Now()
	// for i := num; i < num+50; i++ {
	// 	err := models.GetDBObject().Exec(fmt.Sprintf("insert into cBalance(user_id,status) values(%d,1);", i)).Error
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	// log.Println("insert ", time.Now().Sub(start))

	// start = time.Now()
	// for i := num + 50; i < num+100; i++ {
	// 	balance := models.Balance{UserId: int64(i), Status: 1}
	// 	err := models.GetDBObject().Create(&balance).Error
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	// log.Println("create ", time.Now().Sub(start))

	// start = time.Now()
	// for i := num + 100; i < num+150; i++ {
	// 	err := models.GetDBObject().Create(&models.Balance{UserId: int64(i), Status: 1}).Error
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	// log.Println(" create1 ", time.Now().Sub(start))
	// log.Println(num)
	// start := time.Now()
	// tx := models.GetDBObject().Begin()
	// for i := 10; i < 500; i++ {
	// 	err := tx.Exec(fmt.Sprintf("insert into cBalance(user_id,status) values(%d,1);", i)).Error
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	// tx.Commit()
	// tx.Commit()
	// log.Println("tx insert ", time.Now().Sub(start))
	// start := time.Now()
	// tx1 := models.GetDBObject().Begin()
	// for i := 100; i < 800; i++ {
	// 	err := tx1.Exec(fmt.Sprintf("insert into cBalance(user_id,status) values(%d,1);", i)).Error
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	// tx1.Commit()
	// log.Println("tx insert ", time.Now().Sub(start))

	// start = time.Now()
	// tx1 := models.GetDBObject().Begin()
	// for i := num + 200; i < num+250; i++ {
	// 	err := tx1.Create(&models.Balance{UserId: int64(i), Status: 1}).Error
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	// tx1.Commit()
	// log.Println("tx create ", time.Now().Sub(start))

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
