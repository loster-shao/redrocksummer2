package main

import (
	"github.com/gin-gonic/gin"

	"redrocksummer2/controller"
	"redrocksummer2/model"
	"redrocksummer2/service"
)

func main() {
	model.InitDB()
	service.InitService()//TODO 这个好复杂
	r := gin.Default()//gin初始化
	//路由
	r.GET("/getGoods", controller.SelectGoods)//找货 基本看完。。。有点问题需要等会实验
	r.POST("/order", controller.MakeOrder)//卖货
	r.POST("/add", service.AddGoods)

	r.Run(":8080")
}
//学长tql，随便写的代码都那么牛逼


