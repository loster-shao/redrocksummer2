package service

import (
	"fmt"
	_ "fmt"
	"log"
	"redrocksummer2/model"
	"sync"
	"time"
)

type User struct {
	UserId string
	GoodsId  uint
}

var OrderChan = make(chan User, 1024)

var ItemMap = make(map[uint]*Item)

type Item struct {
	ID        uint   // 商品id
	Name      string // 名字
	Total     int    // 商品总量
	Left      int    // 商品剩余数量
	IsSoldOut bool   // 是否售罄
	leftCh    chan int       //应该是剩余管道先这么命名吧
	sellCh    chan int       //出售管道
	done      chan struct{}  //TODO 这是啥？没看懂。。。。
	Lock      sync.Mutex     //TODO 今天学的，还没完全懂。。。
}

//TODO 写一个定时任务，每天定时从数据库加载数据到Map！！！！！！！！！！！！！！！！

////加物品
//func AddShelve()  {
//	beginTime := time.Now()
//	// 获取第二天时间
//	nextTime := beginTime.Add(time.Hour * 24)
//	// 计算次日零点，即商品上架的时间
//	offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())
//	fmt.Println(offShelveTime)
//	timer := time.NewTimer(offShelveTime.Sub(beginTime))
//	<-timer.C
//
//	var good  []Goods
//	model.DB.Find(&good)
//	fmt.Println(good)
//	var S []Item
//	for i := 0; i < len(good); i++ {
//		s :=  Item{
//			ID:        good[i].ID,
//			Name:      good[i].Name,
//			Total:     good[i].Num,
//			Left:      good[i].Num,
//			IsSoldOut: false,
//		}
//
//		S := append(S, s)
//		ItemMap = S
//	}
//}

func InitMap() {
	beginTime := time.Now()
	// 获取第二天时间
	nextTime := beginTime.Add(time.Hour * 24)
	// 计算次日零点，即商品上架的时间
	offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())
	fmt.Println(offShelveTime)
	timer := time.NewTimer(offShelveTime.Sub(beginTime))
	//	<-timer.C
	for {
		<-timer.C
		for _, i2 := range ItemMap {
			some := model.Goods{
				Name: i2.Name,
				Num:  i2.Left,
			}
			if err := some.AddGoods(); err != nil{
				log.Println(err)
				return
			}
		}
		some := SelectGoods()
		for _, i2 := range some {
			item := &Item{
				ID:        i2.ID,
				Name:      i2.Name,
				Total:     i2.Num,
				Left:      i2.Num,
				IsSoldOut: false,
				leftCh:    make(chan int),
				sellCh:    make(chan int),
			}
			ItemMap[item.ID] = item
		}
		timer.Reset(time.Hour * 24)
	}
}

func initMap() {

	item := &Item{
		ID:        1,
		Name:      "测试",
		Total:     100,
		Left:      100,
		IsSoldOut: false,
		leftCh:    make(chan int),  //管道
		sellCh:    make(chan int),  //管道
	}
	ItemMap[item.ID] = item  //TODO map商品ID等于这个结构体应该是这样的
}

func getItem(itemId uint) *Item{
	return ItemMap[itemId]
}

//订购？应该是
func order() {
	for {
		user := <- OrderChan //从订购管道中接受数据
		item := getItem(user.GoodsId)//TODO 获取商品应该是
		item.SecKilling(user.UserId)//
	}
}

func (item *Item) SecKilling(userId string) {

	item.Lock.Lock()//锁
	defer item.Lock.Unlock()//解锁
	// 等价
	// var lock = make(chan struct{}, 1}
	// lock <- struct{}{}
	// defer func() {
	// 		<- lock
	// }
	if item.IsSoldOut {
		return
	}
	item.BuyGoods(1)

	MakeOrder(userId, item.ID,1)


}


// 定时下架
func (item *Item) OffShelve() {
	beginTime := time.Now()

	// 获取第二天时间
	nextTime := beginTime.Add(time.Hour * 24)
	// 计算次日零点，即商品下架的时间
	offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())

	timer := time.NewTimer(offShelveTime.Sub(beginTime))

	<-timer.C//TODO 这个有何用意？
	delete(ItemMap, item.ID)//删除ID
	close(item.done)

}

// 出售商品
func (item *Item) SalesGoods() {
	for {
		//选择
		select {
		//num由哪个管道发进来就运行哪个
		case num := <-item.sellCh:
			if item.Left -= num; item.Left <= 0 {
				item.IsSoldOut = true
			}

		case item.leftCh <- item.Left:

		case <-item.Done():
			log.Println("我自闭了")
			return
		}
	}
}

func (item *Item) Done() <-chan struct{} {
	//done不知道是啥
	if item.done == nil {
		item.done = make(chan struct{})
	}
	d := item.done
	return d
}

//TODO 监视器？？？啥意思？为啥要取这名字
func (item *Item) Monitor() {
	go item.SalesGoods()
}

// 获取剩余库存
func (item *Item) GetLeft() int {
	var left int
	left = <-item.leftCh
	return left
}

// 购买商品
func (item *Item) BuyGoods(num int) {
	item.sellCh <- num//TODO 从数量管道中传到sellCh管道？？？
}

func InitService() {
	initMap()
	//遍历item切片
	for _,item := range ItemMap{
		item.Monitor()//监视器？
		go item.OffShelve()
		go InitMap()
	}
	for i := 0; i < 10; i++ {
		go order()
	}
}
