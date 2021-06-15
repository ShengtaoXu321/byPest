package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"main.go/model"
	. "main.go/mongo"
	"main.go/token"
	"net/http"
	"strings"
	"time"
)

// 给硬件和前端提供接口
func RouterInit() {
	// 路由监听，给硬件的
	r := gin.Default()
	r.Use(Cors())
	r.POST("/receive", handle1)

	// 网页查询的路由--Post
	r.POST("/history", handle2)
	r.POST("/latest", handle3)
	r.POST("/germs", handle4)

	//r.GET("/latest", handle3)

	r.Run(":18088")

}

// 访问计算器接口获取虫害数据
func GetGermsInit() {
	for true {
		fmt.Println("------")
		// 获取加密的token
		s1, t1 := token.GenerateToken()
		url1 := "https://open-gate.daqiuyin.com/v1"
		body := model.GetGerms{ // 实例化一个请求体
			Method: "GET",
			Path:   "/sc-v2/calculator/60ae5a0dffaf33c74592d6c8",
			Data:   nil,
		}

		//头部信息的封装
		buf := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buf)
		if err := encoder.Encode(body); err != nil {
			log.Println("头部信息编码失败！", err)
		}

		request, err := http.NewRequest(http.MethodPost, url1, buf)
		if err != nil {
			log.Println("头部加载绑定失败！", err)
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-DAQIUYIN-ID", "5f45d17204da596300000002")
		request.Header.Add("X-DAQIUYIN-SIGN", s1)
		request.Header.Add("X-DAQIUYIN-DATE", t1)

		// 解决docker部署时，出现证书错误
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		// 超时
		var client = http.Client{
			Timeout:   20 * time.Second,
			Transport: tr,
		}
		response, err1 := client.Do(request)
		if err1 != nil {
			log.Println("发送POST请求失败！", err)
			return
		}

		defer response.Body.Close()
		//fmt.Println("开始接收数据")
		newBody, _ := ioutil.ReadAll(response.Body)
		str := string(newBody)
		//fmt.Println(string(newBody))

		// 对获取的数据进行解析 -- 使用gjson
		var GermsMap model.ParsingGerms
		data1 := gjson.Get(str, "data").String()
		//fmt.Println(data1)
		data2 := gjson.Get(data1, "algorithm").String()
		data3 := gjson.Get(data2, "calculated").String()
		gemData := gjson.Get(data3, "germs").Int()
		senData := gjson.Get(data2, "sensor_data").String()
		ymd := gjson.Get(senData, "ymd").String()
		T := gjson.Get(senData, "T").Int()
		DC := gjson.Get(senData, "DC").Float()

		// 将解析后的数据存入结构体
		GermsMap.Germs = gemData
		GermsMap.Ymd = ymd + " " + "00:00:00"
		// 时间转时间戳
		loc, _ := time.LoadLocation("Asia/Shanghai")
		tt, _ := time.ParseInLocation("2006-01-02 15:04:05", GermsMap.Ymd, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
		GermsMap.Time = tt.Unix()
		GermsMap.T = T
		GermsMap.Dc = DC
		fmt.Println("拿到的数据信息")
		fmt.Println(GermsMap)

		InsertGerms(GermsMap)

		// 睡眠
		time.Sleep(24 * time.Hour)
	}

}

// 对硬件发送的请求进行判断匹配
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

		fmt.Println("接收到的数据是", recData)

		Insert(recData)
		c.JSON(http.StatusOK, "访问接口成功")

	}

}

// 对接收到的网页请求数据进行解析
func handle2(c *gin.Context) {
	// 1. 将获取到的数据进行校验
	var HistData model.HistoryData
	err := c.ShouldBindJSON(&HistData) // 数据校验--校验格式 [如果数据校验成功，数据存到RecData]
	if err != nil {
		log.Println("网页请求历史数据格式有误！", err)
		c.JSON(http.StatusOK, "网页请求历史数据格式有误！")
		return // 如果忽略return也可以，因为：用户的一次请求只能拥有一次响应
	}
	log.Println("网页请求虫害历史数据格式正确...")

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
		"len":  len(rsp),
		"rsp":  RecData,
	})
}

// 对网页请求的最新数据进行解析判断
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
	// 获取相同时间戳的最新数据
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
		"len":  flag,
		"rsp":  rsp[:100],
	})

}

// 霉变获取数据函数
func handle4(c *gin.Context) {
	// 1. 将获取到的数据进行校验
	var HistGerms model.HistoryData
	err := c.ShouldBindJSON(&HistGerms) // 数据校验--校验格式 [如果数据校验成功，数据存到RecData]
	if err != nil {
		log.Println("网页请求历史数据格式有误！", err)
		c.JSON(http.StatusOK, "网页请求历史数据格式有误！")
		return // 如果忽略return也可以，因为：用户的一次请求只能拥有一次响应
	}
	log.Println("网页请求霉变历史数据格式正确...")

	errCode, rsp := HistoryGerms(HistGerms)
	if errCode != model.SUCESS {
		log.Println(errCode)
		c.JSON(http.StatusOK, gin.H{
			"code": errCode,
			"msg":  "数据查找失败！",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": errCode,
			"len":  len(rsp),
			"rsp":  rsp,
		})
	}
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
