package dao

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database			// 数据库
var personColl *mongo.Collection	// 集合

func Connect() {
	// 设置客户端选项，查询超时时间10秒
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetSocketTimeout(10*time.Second)
	// 连接 MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
	    log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	db = client.Database("local")
	personColl = db.Collection("person")
}

// QueryOne 查找一个文档
func QueryOne(filter interface{}) bson.M {
	var res bson.M
	err := personColl.FindOne(context.TODO(),filter).Decode(&res)
	if err != nil {
		//log.Fatal(err)	// 查不到会报错
		fmt.Println("no match data")
	} else {
		fmt.Println("%+v\n",res)
	}
	return res
}


func QueryMany(filter interface{}) []bson.M {
	var res []bson.M
	cursor, err := personColl.Find(context.TODO(),filter)
	if err != nil {
		//log.Fatal(err)	// 查不到会报错
		fmt.Println("no match data")
	} else {
		defer cursor.Close(context.TODO())
		err = cursor.All(context.TODO(), &res)
		if err != nil {
			fmt.Println(err)
		}
	}
	return res
}