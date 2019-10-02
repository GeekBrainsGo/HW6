package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

// Blog - объект блога
type Blog struct {
	Mongo `inline`
	Title string `bson:"title"`
	Body  string `bson:"article"`
}

func (m *Blog) GetMongoCollectionName() string {
	return "myblogs"
}

// GetAllBlogs - получение всех блогов
func GetAllBlogs(ctx context.Context, db *mongo.Database) ([]Blog, error) {

	blog := Blog{}

	// fmt.Println(blog.GetMongoCollectionName())

	col := db.Collection(blog.GetMongoCollectionName())

	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var blogs []Blog
	if err := cur.All(ctx, &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

// AddBlog - добавляет блог в БД
func (blog *Blog) AddBlog(ctx context.Context, db *mongo.Database) error {

	col := db.Collection(blog.GetMongoCollectionName())
	_, err := col.InsertOne(ctx, blog)
	if err != nil {
		return err
	}
	return nil
}

// Delete - удалят блог из базы
func (blog *Blog) Delete(ctx context.Context, db *mongo.Database) error {

	col := db.Collection(blog.GetMongoCollectionName())
	_, err := col.DeleteOne(ctx, bson.M{"_id": blog.ID})

	return err
}

// GetBlog - получение всех блогов
func GetBlog(ctx context.Context, db *mongo.Database, id primitive.ObjectID) (*Blog, error) {

	blog := Blog{}
	col := db.Collection(blog.GetMongoCollectionName())

	res := col.FindOne(ctx, bson.M{"_id": id})
	if err := res.Decode(&blog); err != nil {
		return nil, err
	}
	return &blog, nil

}

// UpdateBlog - обновляет объект в БД
func (blog *Blog) UpdateBlog(ctx context.Context, db *mongo.Database) error {

	col := db.Collection(blog.GetMongoCollectionName())
	_, err := col.ReplaceOne(ctx, bson.M{"_id": blog.ID}, blog)
	return err
}
