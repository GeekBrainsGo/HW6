package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m *Mongo) GetMongoCollectionName() string {
	panic("GetMongoCollectionName not implemented")
	return ""
}

// BlogItem - объект блога
type BlogItem struct {
	Mongo `inline`
	// Title string `json:"title"`
	// Body  string `json:"article"`
	Title string `bson:"title"`
	Body  string `bson:"article"`
}

// BlogItems - список блогов
type BlogItems []BlogItem

// GetAllBlogItems - получение всех блогов
func GetAllBlogItems(ctx context.Context, db *mongo.Database) ([]BlogItems, error) {

	// blog := BlogItem{}

	//	col := db.Collection(blog.GetMongoCollectionName())
	col := db.Collection("myblogs")

	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	fmt.Println(cur.Current.String())

	blogs := []BlogItems{}
	if err := cur.All(ctx, &blogs); err != nil {
		fmt.Println("Error")

		return nil, err
	}

	return blogs, nil
}

// // AddBlog - обновляет объект в БД
// func (blog *BlogItem) AddBlog(db *gorm.DB) error {

// 	db.Create(&blog)

// 	return nil
// }

// // UpdateBlog - обновляет объект в БД
// func (blog *BlogItem) UpdateBlog(db *gorm.DB) error {

// 	db.Save(&blog)

// 	return nil
// }

// // GetAllBlogItems - получение всех блогов
// func GetBlogItem(db *gorm.DB, id uint) (BlogItem, error) {

// 	blog := BlogItem{}

// 	// Get record with primary key (only works for integer primary key)
// 	db.First(&blog, id)

// 	return blog, nil
// }

// // Delete - удалят объект из базы
// func (blog *BlogItem) Delete(db *gorm.DB) error {

// 	// soft deleted
// 	db.Delete(&blog) // UPDATE blogitems SET deleted_at="2013-10-29 10:23" WHERE id = NUM;
// 	return nil
// }
