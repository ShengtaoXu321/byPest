package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"main.go/model"
	"time"
)

var db_name *mongo.Collection
var db_name1 *mongo.Collection

func MongoInit() {
	// 配置连接
	clientOptions := options.Client().ApplyURI("mongodb://117.78.34.82:18900/")

	// 连接数据库
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("数据库连接失败！", err)
		return
	}
	fmt.Println("mongo数据库连接成功！")
	db := client.Database("ByPest")
	db_name = db.Collection("Pest")
	db_name1 = db.Collection("Germs")
}

// 硬件数据插入数据库Pest集合
func Insert(buf model.PestDate) {
	data := buf.Data
	l := len(data)
	fmt.Println(l)
	for i := 1; i <= l; i++ {
		var dataInDb model.InterData
		dataInDb = data[i-1]
		dataInDb.Time = time.Now().Unix() //改用当前时间戳代替发送数据的时间戳
		_, err := db_name.InsertOne(context.TODO(), dataInDb)
		if err != nil {
			log.Println("数据插入数据库失败！", err)
			continue
		}
		fmt.Println("数据插入成功！", dataInDb.Id)
		continue

	}
}

// 霉变数据插入Germs集合
func InsertGerms(buf model.ParsingGerms) {
	fmt.Println("霉变数据开始插入数据库")
	buf.InsTime = time.Now().Unix() // 霉变加入时间戳
	iResult, err := db_name1.InsertOne(context.TODO(), buf)
	if err != nil {
		log.Println("数据插入失败！", err)
	}
	fmt.Println("霉变数据插入成功！", iResult.InsertedID)

}

// 网页请求的历史数据查询
func History(HistData model.HistoryData) (int, []model.InterData) {
	var res []model.InterData
	//res := []model.InterData{}

	opts := options.Find()
	opts.SetSort(bson.D{{"time", -1}}) // 时间戳从小到大排序，设置可选规则
	ctx := context.Background()        // 全部表格
	filter := bson.M{
		"time": bson.M{"$gte": (HistData.StartTime) / 1000, "$lte": (HistData.EndTime) / 1000},
	}

	// 进行查询逻辑
	if cur, err := db_name.Find(ctx, filter, opts); err != nil {
		log.Println("查询历史数据失败！", err)
		return model.SEARCH_ERR, res
	} else { // 进行数据遍历
		for cur.Next(ctx) {
			var CI model.InterData
			err := cur.Decode(&CI)
			if err != nil {
				log.Println("对历史数据进行解析失败！", err)
				return model.MARSH_ERR, res
			}
			res = append(res, CI)

		}
		if len(res) == 0 {
			log.Println("查询到历史数据为空！")
			return model.SEARCH_NULL, res
		} else {
			return model.SUCESS, res
		}

	}

}

// 网页请求的最新数据查询
func Latest() (int, []model.InterData) {
	var res []model.InterData
	//res := []model.InterData{}

	opts := options.Find()
	opts.SetSort(bson.D{{"time", -1}}) // 时间戳从小到大排序，设置可选规则
	ctx := context.Background()        // 全部表格
	filter := bson.M{}

	// 进行查询逻辑
	if cur, err := db_name.Find(ctx, filter, opts); err != nil {
		log.Println("查询最新数据失败！", err)
		return model.SEARCH_ERR, res
	} else { // 进行数据遍历
		for cur.Next(ctx) {
			var CI model.InterData
			err := cur.Decode(&CI)
			if err != nil {
				log.Println("对最新的数据进行解析失败！", err)
				return model.MARSH_ERR, res
			}
			res = append(res, CI)

		}
		if len(res) == 0 {
			log.Println("查询到最新的数据为空！")
			return model.SEARCH_NULL, res
		} else {
			return model.SUCESS, res
		}

	}

}

// 网页请求霉变历史数据
func HistoryGerms(HistData model.HistoryData) (int, []model.ParsingGerms) {
	var res []model.ParsingGerms
	//res := []model.InterData{}

	opts := options.Find()
	opts.SetSort(bson.D{{"InsTime", -1}}) // 时间戳从小到大排序，设置可选规则
	ctx := context.Background()           // 全部表格
	filter := bson.M{
		"InsTime": bson.M{"$gte": (HistData.StartTime) / 1000, "$lte": (HistData.EndTime) / 1000},
	}

	// 进行查询逻辑
	if cur, err := db_name1.Find(ctx, filter, opts); err != nil {
		log.Println("查询霉变历史数据失败！", err)
		return model.SEARCH_ERR, res
	} else { // 进行数据遍历

		for cur.Next(ctx) {
			var CI model.ParsingGerms
			err := cur.Decode(&CI)
			if err != nil {
				log.Println("对数据进行解析失败！", err)
				return model.MARSH_ERR, res
			}
			res = append(res, CI)

		}
		if len(res) == 0 {
			log.Println("查询到数据为空！")
			return model.SEARCH_NULL, res
		} else {
			return model.SUCESS, res
		}

	}

}
