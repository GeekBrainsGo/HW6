package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	Mongo    `inline`
	Title    string `json:"title,omitempty"`
	Contents string `json:"contents,omitempty"`
}

type Blogs []Blog

func (b *Blog) GetMongoCollectionName() string {
	return "blogs"
}

func GetBlogs(ctx context.Context, db *mongo.Database) (Blogs, error) {
	b := Blog{}
	col := db.Collection(b.GetMongoCollectionName())

	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var blogs Blogs
	if err := cur.All(ctx, &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (b *Blog) Insert(ctx context.Context, db *mongo.Database) (*Blog, error) {
	col := db.Collection(b.GetMongoCollectionName())
	_, err := col.InsertOne(ctx, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func FindBlog(ctx context.Context, db *mongo.Database, id string) (*Blog, error) {
	b := Blog{}
	col := db.Collection(b.GetMongoCollectionName())
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res := col.FindOne(ctx, bson.M{"_id": docID})
	if err := res.Decode(&b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (b *Blog) Update(ctx context.Context, db *mongo.Database) (*Blog, error) {
	col := db.Collection(b.GetMongoCollectionName())
	_, err := col.ReplaceOne(ctx, bson.M{"_id": b.ID}, b)
	return b, err
}

func (b *Blog) Delete(ctx context.Context, db *mongo.Database) (*Blog, error) {
	col := db.Collection(b.GetMongoCollectionName())
	_, err := col.DeleteOne(ctx, bson.M{"_id": b.ID})
	return b, err
}
