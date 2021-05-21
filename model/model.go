package model

// model里面是存放全局的结构体信息


// 1. 定义一个结构体
type PestDate struct {
	Data []InterData `json:"data" binding:"required"`      // 接口体的嵌套
	TotalNums int  `json:"total_nums" binding:"required"`
}

type InterData struct {
	DevType int `json:"devType" binding:"required"`
	Id 	int64	`json:"id" binding:"required"`
	IdDev	string`json:"idDev" binding:"required"`
	PestType int `json:"pestType" binding:"required"`
	Time int64 `json:"time" binding:"required"`
}


