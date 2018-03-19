package models

type RedPacket struct {
	RpId           int64   `gorm:"primary_key" json:"rp_id"`
	Guid           string  `gorm:"not null; default:''; type:varchar(128)" json:"guid"`
	UserId         int64   `gorm:"not null; default:0; type:int; index" json:"user_id"`
	RoomId         int64   `gorm:"not null; default:0; type:int; index" json:"room_id"`
	RedPacketNum   int64   `gorm:"not null; default:0; type:int" json:"red_packet_num"`
	RedPacketMoney float64 `gorm:"not null; default:0; type:decimal(12,2)" json:"red_packet_money"`
	RedPacketType  int64   `gorm:"not null; default:0; type:int" json:"red_packet_type"`
	GetNum         int64   `gorm:"not null; default:0; type:int" json:"get_num"`
	Status         int64   `gorm:"not null; default:0; type:int" json:"status"`
	EndStatus      int64   `gorm:"not null; default:0; type:int" json:"end_status"`
	CreateAt       int64   `gorm:"not null; default:0; type:int; index" json:"create_at"`
	Question       string  `gorm:"not null; type:text" json:"question"`
	Remark1        string  `gorm:"not null; default:''; type:varchar(128); column:remark1" json:"-"`
	Remark2        string  `gorm:"not null; default:''; type:varchar(128); column:remark2" json:"-"`
}

type RedPacketParams struct {
	ID             int64   `gorm:"primary_key" json:"id"`
	RpId           int64   `gorm:"not null; default:0; type:int; index" json:"rp_id"`
	RedPacketNo    int64   `gorm:"not null; default:0; type:int; index" json:"red_packet_no"`
	RedPacketMoney float64 `gorm:"not null; default:0; type:decimal(12,2)" json:"red_packet_money"`
	Status         int64   `gorm:"not null; default:0; type:int" json:"status"`
	Lucky          int64   `gorm:"not null; default:0; type:int" json:"lucky"`
	UserId         int64   `gorm:"not null; default:0; type:int" json:"user_id"`
	UpdateAt       int64   `gorm:"not null; default:0; type:int; index" json:"update_at"`
}

func (RedPacket) TableName () string {
	return "cRedPacket"
}

func (RedPacketParams) TableName () string {
	return "cRedPacketParams"
}

