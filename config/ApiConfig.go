package config

const (
	GEN_RED_PACKET = "genRedPacket" //TODO: 生成红包
)

const (
	RESPONSE_OK    = 0
	RESPONSE_ERROR = 1
	ERROR_MSG      = "系统错误"
)

type GenRedPacket struct {
	WeddingId      int64   `json:"weddingId"`
	RoomId         int64   `json:"roomId"`
	RedPacketMoney float64 `json:"redPacketMoney"`
	RedPacketNum   float64 `json:"redPacketNum"`
	RedPacketType  int64   `json:"redPacketType"`
	Wish           string  `json:"wish"`
	Question       Question
}

type Question struct {
	Title    string   `json:"title"`
	Key      int64    `json:"key"`
	Question []string `json:"question"`
}
