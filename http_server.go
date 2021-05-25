package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"main.go/model"
	"net/http"
	"strings"
)

func RouterInit() {
	// 路由监听，给硬件的
	r := gin.Default()
	r.Use(Cors())
	r.POST("/receive", handle1)

	// 网页查询的路由--Post
	r.POST("/history", handle2)
	r.POST("/latest", handle3)
	//r.GET("/latest", handle3)

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
	var flag int = 0
	for j := 0; j < len(rsp); j++ {
		if rsp[j].Time == rsp[0].Time {
			flag++
			continue
		} else {
			flag = j
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": errCode,
		"rsp":  rsp[:flag],
	})

}

// 解决跨域的CROS
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
