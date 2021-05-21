package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"main.go/model"
	"net/http"
)

func RouterInit() {
	// 路由监听，给硬件的
	r := gin.Default()
	r.POST("/receive", handle1)

	// 网页查询的路由--Post
	r.POST("/history", handle2)
	r.POST("/latest", handle3)

	r.Run(":18088")

}

func handle1(c *gin.Context) {
	{
		// 1. 将获取到的数据进行校验
		var recData model.PestDate
		err := c.ShouldBindJSON(&recData) // 数据校验--校验格式 [如果数据校验成功，数据存到RecData]
		if err != nil {
			log.Println("硬件上传数据格式有误", err)
			c.JSON(http.StatusOK, "硬件上传数据格式有误")
			return // 如果忽略return也可以，因为：用户的一次请求只能拥有一次响应
		}
		log.Println("硬件上传数据格式正确")
		c.JSON(http.StatusOK, "硬件上传数据格式正确")
		//fmt.Println(recData)

		Insert(recData)
		c.JSON(http.StatusOK, "访问接口成功")
	}
}

func handle2(c *gin.Context) {
	// 1. 将获取到的数据进行校验
	var HistData model.HistoryData
	err := c.ShouldBindJSON(&HistData) // 数据校验--校验格式 [如果数据校验成功，数据存到RecData]
	if err != nil {
		log.Println("网页请求历史数据格式有误！", err)
		c.JSON(http.StatusOK, "网页请求历史数据格式有误！")
		return // 如果忽略return也可以，因为：用户的一次请求只能拥有一次响应
	}
	log.Println("网页请求历史数据格式正确...")

	//fmt.Println(recData)
	// 执行数据库操作函数History()

	errCode, rsp := History(HistData)
	if errCode != model.SUCESS {
		log.Println(errCode)
		c.JSON(http.StatusOK, gin.H{
			"code": errCode,
			"msg":  "数据查找失败！",
		})
		return
	}

	//var RecData=map[string]interface{}{}
	var RecData = make(map[string]interface{})
	RecData["historyData"] = rsp
	lMax := len(rsp)
	RecData["total_nums"] = lMax

	c.JSON(http.StatusOK, gin.H{
		"code": errCode,
		"rsp":  RecData,
	})
}

func handle3(c *gin.Context) {
	// 执行数据库操作函数Latest()
	errCode, rsp := Latest()
	if errCode != model.SUCESS {
		log.Println(errCode)
		c.JSON(http.StatusOK, gin.H{
			"code": errCode,
			"msg":  "查找最新数据失败！",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": errCode,
		"rsp":  rsp[0],
	})

}
