package bean

import (
	"github.com/jinzhu/gorm"
)

// Order 下单
type Order struct {
	gorm.Model
	OrderCode string `form:"ordercode" json:"ordercode"` // 订单编码
	Amount    string `form:"amount" json:"amount"`       // 订单金额
	GoodsID   string `form:"goodsid" json:"goodsid"`     // 货物id
	GoodName  string `form:"goodname" json:"goodname"`   // Iphone(第三期）
	BuyTime   string `form:"buytime" json:"buytime"`     // 购买时间
	UserName  string `form:"username" json:"username"`   // 购买用户名称
	Type      int8   `form:"type" json:"type"`           // 类型  0为下单，1为抽奖
	Desc      string `form:"desc" json:"desc"`           // 用户购买的编码
	TxID      string `form:"txid" json:"txid"`           // 入账txid
}
