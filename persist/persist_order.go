package persist

import (
	"github.com/Eric-GreenComb/one-sendeth/bean"
)

// CreateOrder CreateOrder Persist
func (persist *Persist) CreateOrder(order bean.Order) error {
	err := persist.db.Create(&order).Error
	return err
}

// OrderInfo OrderInfo Persist
func (persist *Persist) OrderInfo(orderCode string) (bean.Order, error) {

	var order bean.Order

	err := persist.db.Table("orders").Where("order_code = ?", orderCode).First(&order).Error

	return order, err
}

// GetAllOrders GetAllOrders Persist
func (persist *Persist) GetAllOrders(catid, patchid string) ([]bean.Order, error) {

	var orders []bean.Order

	err := persist.db.Table("orders").Where("goods_id = ? AND good_name = ?", catid, patchid).Find(&orders).Error

	return orders, err
}
