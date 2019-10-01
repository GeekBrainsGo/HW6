package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostItem - объект Post
type PostItem struct {
	Mongo            `inline`
	Idn              string `bson:"idn"`
	Title            string `bson:"title"`
	Date             string `bson:"date"`
	SmallDescription string `bson:"desc"`
	Description      string `bson:"content"`
}

// GetMongoCollectionName - Перегруженный метод возвращающий имя коллекции структуры
func (p *PostItem) GetMongoCollectionName() string {
	return "posts"
}

// PostItemSlice - массив задач
type PostItemSlice []PostItem

// Insert - добавляет пост в БД
func (post *PostItem) Insert(ctx context.Context, db *mongo.Database) (*PostItem, error) {
	col := db.Collection(post.GetMongoCollectionName())
	_, err := col.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// Delete - удалят объект из базы
func (post *PostItem) Delete(ctx context.Context, db *mongo.Database) (*PostItem, error) {
	col := db.Collection(post.GetMongoCollectionName())
	_, err := col.DeleteOne(ctx, bson.M{"idn": post.Idn})
	return post, err
}

// Update - изменяет пост в БД
func (post *PostItem) Update(ctx context.Context, db *mongo.Database) (*PostItem, error) {
	col := db.Collection(post.GetMongoCollectionName())
	_, err := col.ReplaceOne(ctx, bson.M{"idn": post.Idn}, post)
	return post, err
}

// GetAllPosts - получение всех постов
func GetAllPosts(ctx context.Context, db *mongo.Database) (PostItemSlice, error) {
	p := PostItem{}
	col := db.Collection(p.GetMongoCollectionName())

	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	posts := []PostItem{}
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
