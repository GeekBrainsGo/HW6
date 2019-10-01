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

// GetBlog - получение всех блогов
func GetBlog(ctx context.Context, db *mongo.Database, id string) (*Blog, error) {

	// blog := Blog{}

	// // Get record with primary key (only works for integer primary key)
	// db.First(&blog, id)

	// return blog, nil

	blog := Blog{}
	col := db.Collection(blog.GetMongoCollectionName())
	blogID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := col.FindOne(ctx, bson.M{"_id": blogID})
	if err := res.Decode(&blog); err != nil {
		return nil, err
	}
	return &blog, nil

}

// // AddBlog - обновляет объект в БД
// func (blog *Blog) AddBlog(db *gorm.DB) error {

// 	db.Create(&blog)

// 	return nil
// }

// // UpdateBlog - обновляет объект в БД
// func (blog *Blog) UpdateBlog(db *gorm.DB) error {

// 	db.Save(&blog)

// 	return nil
// }

// // Delete - удалят объект из базы
// func (blog *Blog) Delete(db *gorm.DB) error {

// 	// soft deleted
// 	db.Delete(&blog) // UPDATE Blogs SET deleted_at="2013-10-29 10:23" WHERE id = NUM;
// 	return nil
// }
