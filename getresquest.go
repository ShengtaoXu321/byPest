package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type GetGerms struct {
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Data   interface{} `json:"data"`
}

type ParsingGerms struct {
	Germs   int64   `bson:"germs" json:"germs"`
	Ymd     string  `bson:"Ymd" json:"Ymd"`
	T       int64   `bson:"T" json:"T"`
	Dc      float64 `bson:"DC" json:"DC"`
	Time    int64   `bson:"time" json:"Time"`
	InsTime int64   `bson:"InsTime" json:"InsTime"`
}

func main() {
	GetGermsInit()
}

func GetGermsInit() {
	for true {
		fmt.Println("------")
		// 获取加密的token
		s1, t1 := GenerateToken()
		url1 := "https://open-gate.daqiuyin.com/v1"
		body := GetGerms{ // 实例化一个请求体
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

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		// 超时
		var client = http.Client{
			Transport: tr,
			Timeout: 20 * time.Second,
		}
		response, err1 := client.Do(request)
		if err1 != nil {
			log.Println("发送POST请求失败！", err1)
			return
		}

		defer response.Body.Close()
		fmt.Println("开始接收数据")
		newBody, _ := ioutil.ReadAll(response.Body)
		str := string(newBody)
		fmt.Println(string(newBody))

		// 对获取的数据进行解析 -- 使用gjson
		fmt.Println("开始解析数据")
		var GermsMap ParsingGerms
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
		loc, err2 := time.LoadLocation("Asia/Shanghai")
		if err2!=nil{
			log.Println(err2)
		}
		fmt.Println(loc)
		tt, err3:= time.ParseInLocation("2006-01-02 15:04:05", GermsMap.Ymd, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
		if err3!=nil{
			log.Println(err3)
		}
		GermsMap.Time = tt.Unix()
		GermsMap.T = T
		GermsMap.Dc = DC
		fmt.Println("拿到的数据信息")
		fmt.Println(GermsMap)

		// 睡眠
		time.Sleep(5 * time.Minute)
	}

}

// 这是用来生成与获取计算器详情接口需要的字段
func GenerateToken() (s string, t string) {
	key := "6HlROpEeOhlnZ5TKzDxMm7Hdaxw+jWbPXv1LDQj/B3="
	//  很重要：时间转换， int64不能直接转换成string
	t = strconv.FormatInt((time.Now().UnixNano() / 1e6), 10)

	//fmt.Println(t)
	path := "GET&/sc-v2/calculator/60ae5a0dffaf33c74592d6c8&" + t
	hmac := hmac.New(sha256.New, []byte(key))
	hmac.Write([]byte(path))
	s = base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	//fmt.Println(base64.StdEncoding.EncodeToString(hmac.Sum(nil)))
	return s, t
}


// Post "https://open-gate.daqiuyin.com/v1": x509: certificate signed by unknown authority