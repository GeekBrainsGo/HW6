/*
 * HomeWork-6: Mongo in BeeGo
 * Created on 28.09.19 22:16
 * Copyright (c) 2019 - Eugene Klimov
 */

package models

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"time"
)

// MDB & Logger globals (it is normal for BeeGo)
var (
	MDB *mongo.Client
	Lg  *logs.BeeLogger
)

// Post is the base post type.
type Post struct {
	OID     primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ID      string             `json:"-" bson:"-"`
	Title   string             `json:"title"`
	Date    time.Time          `json:"-" bson:"updated_at"`
	Summary string             `json:"summary"`
	Body    template.HTML      `json:"body"`
	Created time.Time          `json:"-" bson:"created_at"`
	Deleted time.Time          `json:"-" bson:"deleted_at"`
}

//DBPosts is type dbPosts map[string]Post
type DBPosts struct {
	Collection *mongo.Collection
	ctx        context.Context
	Posts      []Post
	Lg         *logs.BeeLogger
	Error
}

// NewPosts creates new DBPosts with DB link
func NewPosts() *DBPosts {
	dbName := beego.AppConfig.String("DBNAME")
	col := MDB.Database(dbName).Collection(Post{}.TableName())
	return &DBPosts{
		Collection: col,
		ctx:        context.TODO(),
		Lg:         Lg,
		Error:      Error{Lg: Lg},
	}
}

// TableName returns name for post's table (need for ORM & Mongo).
func (Post) TableName() string {
	table := beego.AppConfig.String("TABLENAME")
	if table == "" {
		return "posts"
	}
	return table
}

// Date2Norm normalizes date to local format for view in browsers.
func (p *Post) Date2Norm() string {
	dt := beego.AppConfig.String("DATETIME")
	if dt == "" {
		dt = "02.01.2006 15:04:05"
	}
	s, off := time.Now().Zone()
	return p.Date.Add(time.Second * time.Duration(off)).Format(fmt.Sprintf("%s %s", dt, s))
}

// GetPosts gets one or all posts.
func (d *DBPosts) GetPosts(id string) error {
	post := Post{}
	if id == "" { // all posts
		opts := options.Find()
		opts.SetSort(bson.D{{"updated_at", -1}})
		cur, err := d.Collection.Find(d.ctx, bson.M{"deleted_at": time.Unix(0, 0)}, opts)
		if err != nil {
			return fmt.Errorf("error find all posts: %v", err)
		}
		err = cur.All(d.ctx, &d.Posts)
		if err != nil {
			return fmt.Errorf("error fill post's slice from find results: %v", err)
		}
		for i := range d.Posts {
			d.Posts[i].ID = d.Posts[i].OID.Hex()
		}
	} else { // one post
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("error converting post ID to objectID: %v", err)
		}
		err = d.Collection.FindOne(d.ctx, bson.M{
			"_id": objID, "deleted_at": time.Unix(0, 0),
		}).Decode(&post)
		if err != nil {
			return fmt.Errorf("post not found: %s", id)
		}
		d.Posts = append(d.Posts, post)
	}
	return nil
}

// CreatePost creates post.
func (d *DBPosts) CreatePost() error {
	//d.Posts[0].OID = primitive.NewObjectID() // or omitempty in Post
	d.Posts[0].Date = time.Now()
	d.Posts[0].Created = time.Now()
	d.Posts[0].Deleted = time.Unix(0, 0)
	_, err := d.Collection.InsertOne(d.ctx, d.Posts[0])
	if err != nil {
		return fmt.Errorf("error insert one post: %v", err)
	}
	return nil
}

// DeletePost deletes one post.
func (d *DBPosts) DeletePost(id string) error {
	return d.UpdatePost(id, true)
}

// UpdatePost updates post.
func (d *DBPosts) UpdatePost(id string, isDelete bool) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error converting post ID to objectID: %v", err)
	}
	var upd bson.D
	if isDelete {
		upd = bson.D{
			{"$set", bson.D{
				{"deleted_at", time.Now()},
			}},
		}
	} else {
		upd = bson.D{
			{"$set", bson.D{
				{"title", d.Posts[0].Title},
				{"summary", d.Posts[0].Summary},
				{"body", d.Posts[0].Body},
				{"updated_at", time.Now()},
			}},
		}
	}
	res, err := d.Collection.UpdateOne(d.ctx, bson.M{"_id": objID}, upd)
	if err != nil {
		return fmt.Errorf("error update post: %v", err)
	}
	if res.ModifiedCount == 0 {
		d.Lg.Warning("post not found: %s", id)
	}
	return nil
}
