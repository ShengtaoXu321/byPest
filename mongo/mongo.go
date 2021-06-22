package mongo

import (
	"context"
	"data/model"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"strconv"
	"strings"
)
var db_name *mongo.Collection
var db_name1 *mongo.Collection


func MongoInit()  {
	// 配置连接
	clientOption:=options.Client().ApplyURI("mongodb://117.78.34.82:18900/")

	// 开始连接数据库
	client,err:=mongo.Connect(context.TODO(),clientOption)
	if err!=nil{
		log.Println("数据库连接失败！",err)
		return
	}
	fmt.Println("mongodb数据库连接成功！")
	// 数据库名称
	db:=client.Database("ByPest")
	db_name=db.Collection("pest1")
	db_name1=db.Collection("Germs")

}


var buf model.PestData


// 数据库插入数据
func HandInsert() {
	buf.DevType=0
	buf.IdDev=genDev()
	buf.PestType=int64(rand.Intn(3))
	timebuf:=[]int64{1622520000,1622606400,1622692800,1622779200,1622865600,1622952000,
		1623038400,1623124800,1623211200,1623297600,1623384000,1623470400,1623556800,1623643200,1623729600,1623816000,
		1623902400,1623988800,1624075200,1624161600,1624248000}
	buf.Time=timebuf[1]
	fmt.Println(buf)
	fmt.Println("开始进行数据库插入操作！")
	iResult,err:=db_name.InsertOne(context.TODO(),buf)
	if err!=nil{
		log.Println("数据库插入失败！",err)
	}
	fmt.Println("数据库插入成功！",iResult.InsertedID)

}

// 随机生成iddev函数
func genDev() string{
	rawData:=RandInt64(769,782)
	hexData:=strconv.FormatInt(rawData,16)
	hexData="0"+ hexData
	HexData:=strings.ToUpper(hexData)
	return HexData
}



// 在线随机区间数生成
func RandInt64(min, max int64) int64 {
	if min>=max ||min==0|| max==0{
		return max
	}
	return rand.Int63n(max-min)+min
}






