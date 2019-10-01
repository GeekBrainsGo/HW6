package models

import (
	"context"
	"html/template"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Post stands for post object.
type Post struct {
	Mongo   `bson:"inline"`
	Hex     string        `json:"hex" bson:"hex"`
	Title   string        `bson:"title" json:"title"`
	Author  string        `bson:"author" json:"author"`
	Content template.HTML `bson:"content" json:"content"`
}

// Posts stands for array of posts.
type Posts []Post

const collection = "posts"

// Insert post to database.
func (p *Post) Insert(db *mongo.Database) error {
	p.OID = primitive.NewObjectID()
	p.Hex = p.OID.Hex()

	col := db.Collection(collection)
	_, err := col.InsertOne(context.TODO(), p)
	if err != nil {
		return err
	}
	return nil
}

// Update updates post in database.
func (p *Post) Update(db *mongo.Database) (*Post, error) {
	col := db.Collection(collection)
	_, err := col.ReplaceOne(context.TODO(), bson.M{"hex": p.Hex}, p)
	return p, err
}

// Delete deletes post from database.
func (p *Post) Delete(db *mongo.Database) (*Post, error) {
	col := db.Collection(collection)
	_, err := col.DeleteOne(context.TODO(), bson.M{"hex": p.Hex})
	return p, err
}

// AllPosts return all posts from database.
func AllPosts(db *mongo.Database) ([]Post, error) {
	col := db.Collection(collection)

	cur, err := col.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	if err := cur.All(context.TODO(), &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
