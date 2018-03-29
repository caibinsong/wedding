package config

//通信返回信息结构体
type Response struct {
	Code int64       `json:"code"`           //状态 0 成功，非0失败
	Msg  string      `json:"msg,omitempty"`  //
	Data interface{} `json:"data,omitempty"` //返回内容
}

//获取发红包人信息与抢红包人的信息
type Res_GetRedPacketInfo struct {
	RpId       int64     `json:"rp_id"`
	Red_Info   RedInfo   `json:"red_info"`
	Other_Info OtherInfo `json:"other_info"`
}

//发红包人信息
type RedInfo struct {
	UserId         int     `json:"rp_id"`
	NickName       string  `json:"nick_name"`
	Pic            string  `json:"pic"`
	Wish           string  `json:"wish"`
	RedPacketNum   int64   `json:"red_packet_num"`
	RedPacketMoney float64 `json:"red_packet_money"`
	RedPacketType  int64   `json:"red_packet_type"`
	Status         int64   `json:"status"`
}

//抢红包人结构体
type OtherInfo struct {
	Count int         `json:"count"`
	List  []OtherList `json:"list"`
}

//抢红包人信息
type OtherList struct {
	UserId   int64   `json:"user_id"`
	NickName string  `json:"nick_name"`
	Pic      string  `json:"pic"`
	UpdateAt int64   `json:"update_at"`
	Money    float64 `json:"money"`
	Lucky    int64   `json:"lucky"`
}

//微信用户列表信息
type GetWXUserInfo struct {
	ActionName string `json:"action_name"`
	Data       string `json:"data"`
}

//微信用户列表信息
type GetWXUserList struct {
	ActionName string `json:"action_name"`
	Data       struct {
		UserList []int `json:"user_list"`
	} `json:"data"`
}

//微信返回信息信息
type WXUserInfoResponse struct {
	Code int        `json:"code"`
	Data WXUserInfo `json:"data"`
	Msg  string     `json:"msg"`
}

//用户返回列表信息
type WXUserListResponse struct {
	Code int          `json:"code"`
	Data []WXUserInfo `json:"data"`
	Msg  string       `json:"msg"`
}

//用户返回信息
type WXUserInfo struct {
	Id            string `json:"id"`
	AppId         string `json:"app_id"`
	NickName      string `json:"nickName"`
	Pic           string `json:"pic"`
	OpenId        string `json:"open_id"`
	Language      string `json:"language"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	LastVisitTime string `json:"last_visit_time"`
	RealName      string `json:"real_name"`
	Phone         string `json:"phone"`
}

//微信回调结构体
type WXPayNotifyReq struct {
	Return_code          string `xml:"return_code"`
	Return_msg           string `xml:"return_msg"`
	Appid                string `xml:"appid"`
	Mch_id               string `xml:"mch_id"`
	Device_info          string `xml:"device_info"`
	Nonce                string `xml:"nonce_str"`
	Sign                 string `xml:"sign"`
	Sign_type            string `xml:"sign_type"`
	Result_code          string `xml:"result_code"`
	Err_code             string `xml:"err_code"`
	Err_code_des         string `xml:"err_code_des"`
	Openid               string `xml:"openid"`
	Is_subscribe         string `xml:"is_subscribe"`
	Trade_type           string `xml:"trade_type"`
	Bank_type            string `xml:"bank_type"`
	Total_fee            int    `xml:"total_fee"`
	Settlement_total_fee int    `xml:"settlement_total_fee"`
	Fee_type             string `xml:"fee_type"`
	Cash_fee             int    `xml:"cash_fee"`
	Cash_fee_Type        string `xml:"cash_fee_type"`
	Transaction_id       string `xml:"transaction_id"`
	Out_trade_no         string `xml:"out_trade_no"`
	Attach               string `xml:"attach"`
	Time_end             string `xml:"time_end"`
}

//微信通信返回信息结构体
type WXPayNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}
