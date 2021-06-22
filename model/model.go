package model

type PestData struct {
	DevType   int64 `bson:"devType" json:"devType"`
	IdDev     string`bson:"idDev" json:"idDev"`
	PestType  int64`bson:"pestType" json:"pestType"`
	Time      int64`bson:"time" json:"time"`
}

// 历史数据结构体
type HistoryData struct {
	StartTime int64 `json:"startTime" binding:"required"`
	EndTime   int64 `json:"endTime" binding:"required"`
}

// 3. 定义获取霉变数据的请求头数据
type GetGerms struct {
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Data   interface{} `json:"data"`
}

// 4. 定义霉变解析后需要的数据
type ParsingGerms struct {
	Germs   int64   `bson:"germs" json:"germs"`
	Ymd     string  `bson:"Ymd" json:"Ymd"`
	T       int64   `bson:"T" json:"T"`
	Dc      float64 `bson:"DC" json:"DC"`
	Time    int64   `bson:"time" json:"Time"`
	Level   int64   `bson:"level" json:"Level"`
	InsTime int64   `bson:"InsTime" json:"InsTime"`
}

// 常量错误码
const (
	SEARCH_ERR  = -2 // 数据库查询错误
	MARSH_ERR   = -1 // 数据库返回的数据解析错误
	SEARCH_NULL = 0  // 数据库没有符合的记录
	SUCESS      = 1  // 数据成功
)
