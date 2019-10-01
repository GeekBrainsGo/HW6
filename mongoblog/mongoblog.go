package main

/*
	Basics Go.
	Rishat Ishbulatov, dated Sep 30, 2019.
	Translate your blog to MongoDB.
*/

import (
	"context"
	"flag"
	"HW6/mongoblog/server"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	flagRootDir := flag.String("rootdir", "./www", "root dir of the server")
	flagServAddr := flag.String("addr", "localhost:8080", "server address")
	flag.Parse()

	lg := NewLogger()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		lg.WithError(err).Fatal("can't get new client")
	}

	err = client.Connect(context.TODO())
	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		lg.WithError(err).Fatal("can't ping to db")
	}

	db := client.Database("mongoblog")
	defer client.Disconnect(context.TODO())

	serv := server.New(lg, *flagRootDir, db)

	go func() {
		err := serv.Start(*flagServAddr)
		if err != nil {
			lg.WithError(err).Fatal("can't run the server")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
}

// NewLogger creates new logger.
func NewLogger() *logrus.Logger {
	lg := logrus.New()
	lg.SetReportCaller(false)
	lg.SetFormatter(&logrus.TextFormatter{})
	lg.SetLevel(logrus.DebugLevel)
	return lg
}
