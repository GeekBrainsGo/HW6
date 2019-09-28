package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Post - Модель поста
type Post struct {
	Mongo   `inline`
	Title   string `bson:"title"`
	Desc    string `bson:"desc"`
	Content string `bson:"content"`
}

// GetMongoCollectionName - Перегруженный метод возвращающий имя коллекции структуры
func (p *Post) GetMongoCollectionName() string {
	return "posts"
}

// Insert - Добавляет объект в базу
func (p *Post) Insert(ctx context.Context, db *mongo.Database) (*Post, error) {
	col := db.Collection(p.GetMongoCollectionName())
	_, err := col.InsertOne(ctx, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetPost - Возвращает пост по его id
func GetPost(ctx context.Context, db *mongo.Database, id string) (*Post, error) {
	p := Post{}
	col := db.Collection(p.GetMongoCollectionName())
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res := col.FindOne(ctx, bson.M{"_id": docID})
	if err := res.Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Update - Обновляет весь пост
func (p *Post) Update(ctx context.Context, db *mongo.Database) (*Post, error) {
	col := db.Collection(p.GetMongoCollectionName())
	_, err := col.ReplaceOne(ctx, bson.M{"_id": p.ID}, p)
	return p, err
}

// Delete - Удаляет пост
func (p *Post) Delete(ctx context.Context, db *mongo.Database) (*Post, error) {
	col := db.Collection(p.GetMongoCollectionName())
	_, err := col.DeleteOne(ctx, bson.M{"_id": p.ID})
	return p, err
}

// GetPosts - Возвращает все посты
func GetPosts(ctx context.Context, db *mongo.Database) ([]Post, error) {
	p := Post{}
	col := db.Collection(p.GetMongoCollectionName())

	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// Find - Находит все поста у которых поле равно значению
func Find(ctx context.Context, db *mongo.Database, field string, value interface{}) ([]Post, error) {
	p := Post{}
	col := db.Collection(p.GetMongoCollectionName())

	cur, err := col.Find(ctx, bson.M{field: value})
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
