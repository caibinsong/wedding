package config

const (
	ServerName = "RedPacketLogic"
)

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

const (
	DBName = "cAuth" //wedding_card
	USER   = "root"         //root
	PASS   = "root"         //tkC42cwy2U3SQwHw
	HOST   = "127.0.0.1"    //172.17.0.5
	DEBUG  = "dev"          //prod
)
