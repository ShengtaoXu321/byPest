package main

import (
	. "main.go/http"
	. "main.go/mongo"
)

// 功能叙述： 接收到JSON数据，对数据进行解析，
//以MQTT的形式传输至云平台架构的下一步---APLUS_MQ

func main() {
	//初始化Mongo数据库

	MongoInit()
	//// 打开路由监听
	go RouterInit()

	// 获取计算器详情数据
	GetGermsInit()

	select {}
}
