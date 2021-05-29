package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

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
