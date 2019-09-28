package main

import (
	"context"
	"fmt"
	"log"
	"serv/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	_ = client.Connect(ctx)
	db := client.Database("geekbrains")

	post := &models.Post{
		Title:   "TITLE",
		Desc:    "456",
		Content: "789",
	}

	post, err := post.Insert(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*post)

	post, err = models.GetPost(ctx, db, "5d8fb379deb49eed85319707")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*post)

	post, err = post.Update(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*post)

	post, err = post.Delete(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*post)

	posts, err := models.GetPosts(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(posts)

}
