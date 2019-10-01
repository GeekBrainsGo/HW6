package main

import (
	"blogMongo/server"
	"context"
	"flag"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	var client *mongo.Client

	flagRootDir := flag.String("rootdir", "./web", "root dir of the server")
	flagServAddr := flag.String("addr", "localhost:8080", "server address")
	// flagMongoDB := flag.String("mongo", "mongodb://localhost:27017", "MongoDB connection string, format: mongodb://host:port")
	// // flagMySQL := flag.String("sql", "root:root@/MyBlogs?parseTime=true", "MySQL connection string, format: user:password@tcp(host:port)/database")
	flag.Parse()

	lg := NewLogger()

	ctx := context.Background()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	// client, err := mongo.NewClient(options.Client().ApplyURI(flagMongoDB))
	if err != nil {
		lg.Fatal(err)
	}
	err = client.Connect(ctx)
	// err = client.Connect(context.TODO())
	if err != nil {
		lg.Fatal(err)
	}

	db := client.Database("myblog")

	serv := server.New(ctx, lg, *flagRootDir, db)
	// serv := server.New(lg, *flagRootDir, db)

	go func() {
		err := serv.Start(*flagServAddr)
		if err != nil {
			lg.WithError(err).Fatal("can't run the server")
		}
	}()

	stopSig := make(chan os.Signal)
	signal.Notify(stopSig, os.Interrupt, os.Kill)
	<-stopSig
	lg.Info("Stop server!")

}
