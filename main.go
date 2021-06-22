package main

import "data/mongo"

func main() {
	mongo.MongoInit()
	for i:=0;i<150;i++{
		mongo.HandInsert()
	}
	mongo.HandInsert()

}


