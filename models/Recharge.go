package models

type Recharge struct {
	Id           int64   `gorm:"primary_key" json:"id"`
	UserId       int64   `gorm:"not null; default:0;  type:int" json:"user_id"`                    //婚礼ID
	Price        float64 `gorm:"not null; default:0;  type:decimal(12,2)" json:"price"`            //打赏用户ID
	PayMoney     float64 `gorm:"not null; default:0;  type:decimal(12,2)" index" json:"pay_money"` //婚礼创建人ID
	StatusResult int64   `gorm:"not null; default:0;  type:int" index" json:"status_result"`       //操作类型
	PayType      int64   `gorm:"not null; default:0;  type:int" json:"pay_type"`                   //操作金额
	PayBill      string  `gorm:"null;  type:varchar(100);" json:"pay_bill"`                        //打赏时间
	OpUniqueNo   string  `gorm:"not null; default:'';  type:varchar(100)" json:"op_unique_no"`     //状态
	OtherPayBill string  `gorm:"not null; default:''; type:varchar(100); json:"other_pay_bill"`    //唯一操作号
	PayTimes     int64   `gorm:"not null; default:0; type:int; json:"pay_times"`                   //创建时间戳
	FailReason   string  `gorm:"not null; default:''; type:varchar(256); json:"fail_reason"`       //备注1
	CreateAt     int64   `gorm:"not null; default:0; type:int; json:"create_at"`                   //备注2
	Remark1      string  `gorm:"not null; default:''; type:varchar(128); column:remark1" json:"-"` //备注2
	Remark2      string  `gorm:"not null; default:''; type:varchar(128); column:remark2" json:"-"` //备注2
}

func (Recharge) TableName() string {
	return "cRecharge"
}
