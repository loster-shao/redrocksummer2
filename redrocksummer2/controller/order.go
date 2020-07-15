package controller

import (
	"github.com/gin-gonic/gin"
	"redrocksummer2/service"
	"strconv"

)


func MakeOrder(ctx *gin.Context) {
	//接受Postman参数
	userId := ctx.PostForm("userId")
	goodsId := ctx.PostForm("goodsId")

	//string->int
	itemId,_ := strconv.Atoi(goodsId)

	//接受User来的数据
	service.OrderChan <- service.User{
		UserId:  userId,
		GoodsId: uint(itemId),
	}
	//TODO 应该是购买成功，或者是订单成功，目前不太清楚
	ctx.JSON(200, gin.H{
		"status": 200,
		"info": "success",
	})
}



