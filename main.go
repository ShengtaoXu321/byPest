package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"log"
	"main.go/model"
	"net/http"
)

// 功能叙述： 接收到JSON数据，对数据进行解析，
//以MQTT的形式传输至云平台架构的下一步---APLUS_MQ

func main() {
	r := gin.Default()
	r.POST("/receive", handle1)
	r.Run(":18088")

}

func handle1(c *gin.Context) {
	{
		// 1. 将获取到的数据进行校验
		var recData model.PestDate
		err:=c.ShouldBindJSON(&recData)   // 数据校验--校验格式 [如果数据校验成功，数据存到RecData]
		if err!=nil{
			log.Println("数据格式有误",err)
			c.JSON(http.StatusOK,"数据格式有误")
			return   // 如果忽略return也可以，因为：用户的一次请求只能拥有一次响应
		}
		log.Println("数据格式正确")
		c.JSON(http.StatusOK,"数据格式正确")

		// 2. 数据以MQTT发送 -- 数据已经存到recData里面
		// 第一步：序列化
		repData,err2:=json.Marshal(recData)
		if err2!=nil{
			log.Println("序列化失败",err2)
		}
		// 2. 定义一个token
		topicMqtt:="by/Pest/"+string(recData.Data[0].Time)

		// 3. MQTT连接
		client:=mqtt.NewClientOptions().AddBroker("117.78.34.82:18211")
		c1:=mqtt.NewClient(client)
		if token:=c1.Connect();token.Wait()&&token.Error()!=nil{
			log.Println("连接MQTT服务器失败：",token.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":20001,
				"message":"连接MQTT服务器失败！",
			})
		}







		//buf := make([]byte, 10240000)
		//n, _ := c.Request.Body.Read(buf) //返回的是： 1. buf的长度 2. err
		//fmt.Println("接收的内容为:", string(buf[:n]))
		//c.JSON(http.StatusOK, "访问接口成功")
	}
}

