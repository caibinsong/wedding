package models

type Spending struct {
	Id          int64   `gorm:"primary_key" json:"id"`
	WeddingId   int64   `gorm:"not null; default:0;  type:int" json:"wedding_id"`                 //婚礼ID
	UserId      int64   `gorm:"not null; default:0;  type:int" json:"user_id"`                    //打赏用户ID
	FounderId   int64   `gorm:"not null; default:0;  type:int" index" json:"founder_id"`          //婚礼创建人ID
	OperateType int64   `gorm:"not null; default:0;  type:int" index" json:"operate_type"`        //操作类型
	Money       float64 `gorm:"not null; default:0;  type:decimal(12,2)" json:"money"`            //操作金额
	CreateTime  string  `gorm:"null;  type:datetime;" json:"create_time"`                         //打赏时间
	Status      int64   `gorm:"not null; default:0;  type:int" json:"status"`                     //状态
	OpUniqueNo  string  `gorm:"not null; default:''; type:varchar(100); json:"op_unique_no"`      //唯一操作号
	CreateAt    int64   `gorm:"not null; default:0; type:int; json:"create_at"`                   //创建时间戳
	Remark1     string  `gorm:"not null; default:''; type:varchar(128); column:remark1" json:"-"` //备注1
	Remark2     string  `gorm:"not null; default:''; type:varchar(128); column:remark2" json:"-"` //备注2
}

func (Spending) TableName() string {
	return "cSpending"
}
