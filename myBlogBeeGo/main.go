/*
 * HomeWork-6: Mongo in BeeGo
 * Created on 28.09.19 22:18
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"myBlog/conf"
	"myBlog/models"
	_ "myBlog/routers"
)

func main() {
	// connect to Mongo
	mdb, err := mongo.NewClient(options.Client().ApplyURI(conf.GetURI()))
	if err != nil {
		log.Fatalln("Can't create MongoDB client:", err)
	}
	models.MDB = mdb
	if err = models.MDB.Connect(context.TODO()); err != nil {
		log.Fatalln("Can't connect to MongoDB server:", err)
	}
	if err = models.MDB.Ping(context.TODO(), nil); err != nil {
		log.Fatalln("Can't ping MongoDB server:", err)
	}

	// set logger
	dbName := beego.AppConfig.String("DBNAME")
	models.Lg = logs.NewLogger(10)
	models.Lg.SetPrefix(fmt.Sprintf("[%s]", dbName))
	models.Lg.Info("Connected to MongoDB")
	defer models.Lg.Close()

	beego.Run()

	if err = models.MDB.Disconnect(context.TODO()); err != nil {
		log.Fatalln("error disconnect from MongoDB", err)
	}
}
