package models

type Client struct {
	Id         int64  `gorm:"primary_key" json:"id"`
	UserId     int64  `gorm:"not null; default:0;  type:int" json:"user_id"`                    //用户ID
	RealName   string `gorm:"not null; default:'';  type:varchar(128)" json:"real_name"`        //真实姓名
	Phone      string `gorm:"not null; default:'';  type:varchar(128)" json:"phone"`            //手机号
	CreateTime string `gorm:"null;  type:datetime;" json:"create_time"`                         //创建时间
	Remark1    string `gorm:"not null; default:''; type:varchar(128); column:remark1" json:"-"` //备注1
	Remark2    string `gorm:"not null; default:''; type:varchar(128); column:remark2" json:"-"` //备注2
}

func (Client) TableName() string {
	return "cClient"
}
