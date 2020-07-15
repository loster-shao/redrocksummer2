package model

import "github.com/jinzhu/gorm"

//订单表单数据
type Order struct {
	gorm.Model
	UserID  string
	GoodsID uint
	Num     int
}

// 下单
func (order *Order)MakeOrder() error{
	return DB.Create(&order).Error//返回创建订单
}


// 查询订单
func GetOrderByUserID(userId string) (orders []Order, err error){
	err = DB.Table("orders").Where("user_id = ?",userId).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders,nil
}