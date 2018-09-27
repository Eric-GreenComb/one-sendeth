package bean

import ()

// OneOrder OneOrder
type OneOrder struct {
	OrderCode string `form:"ordercode" json:"ordercode"` // 订单编码
	Amount    string `form:"amount" json:"amount"`       // 订单总金额 不带单位
	GoodsID   string `form:"goodsid" json:"goodsid"`     // 产品ID
	GoodName  string `form:"goodname" json:"goodname"`   // 产品名称
	BuyTime   string `form:"buytime" json:"buytime"`     // 购买时间
	WinTime   string `form:"wintime" json:"wintime"`     // 开奖时间
	UserName  string `form:"username" json:"username"`   // 购买用户名称
	Desc      string `form:"desc" json:"desc"`           // 描述
	Type      int8   `form:"type" json:"type"`           // 类型  0为下单，1为抽奖
}

// Callback Callback
type Callback struct {
	OrderCode string `form:"ordercode" json:"ordercode"` // 订单编码
	TxID      string `form:"txid" json:"txid"`           // TxID
}
