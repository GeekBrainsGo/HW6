/*
 * HomeWork-6: Mongo in BeeGo
 * Created on 28.09.19 22:17
 * Copyright (c) 2019 - Eugene Klimov
 */

package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/russross/blackfriday/v2"
	"html/template"
	"myBlog/models"
	"net/http"
)

// MainController is.
type MainController struct {
	beego.Controller
}

// GetAllPosts shows all posts in main page.
func (c *MainController) GetAllPosts() {
	posts := models.NewPosts()
	if err := posts.GetPosts(""); err != nil {
		posts.Lg.Error("error get all posts: %s", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Abort("500")
		return
	}
	c.Data["BlogName"] = beego.AppConfig.String("BLOGNAME")
	c.Data["Posts"] = &posts.Posts
	c.TplName = "index.tpl"
}

// GetOnePost shows one posts with full content.
func (c *MainController) GetOnePost() {
	postNum := c.Ctx.Request.URL.Query().Get("id")
	if postNum == "" {
		c.Redirect("/", http.StatusMovedPermanently)
	}
	posts := models.NewPosts()
	if err := posts.GetPosts(postNum); err != nil {
		posts.Lg.Error("error get one post: %s", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Abort("404")
	}
	c.Data["BlogName"] = beego.AppConfig.String("BLOGNAME")
	posts.Posts[0].Body = template.HTML(blackfriday.Run([]byte(posts.Posts[0].Body)))
	posts.Posts[0].ID = posts.Posts[0].OID.Hex()
	c.Data["Post"] = &posts.Posts[0]
	c.TplName = "post.tpl"
}

// DeletePost removes post from DB.
func (c *MainController) DeletePost() {
	postNum := c.Ctx.Input.Param(":id")
	posts := models.NewPosts()
	if err := posts.DeletePost(postNum); err != nil {
		posts.Lg.Error("error delete post: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while delete post")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
}

// GetEditPost shows edit form for edit post.
func (c *MainController) GetEditPost() {
	postNum := c.Ctx.Request.URL.Query().Get("id")
	if postNum == "" {
		c.Redirect("/", http.StatusMovedPermanently)
	}
	posts := models.NewPosts()
	if err := posts.GetPosts(postNum); err != nil {
		posts.Lg.Error("error get one post for edit: %s", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Abort("404")
	}
	c.Data["BlogName"] = beego.AppConfig.String("BLOGNAME")
	posts.Posts[0].ID = posts.Posts[0].OID.Hex()
	c.Data["Post"] = &posts.Posts[0]
	c.TplName = "edit.tpl"
}

// UpdatePost updates post in DB.
func (c *MainController) UpdatePost() {
	postNum := c.Ctx.Input.Param(":id")
	posts := models.NewPosts()
	post, err := c.decodePost()
	if err != nil {
		posts.Lg.Error("error while decoding post body: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while decoding post body")
		return
	}
	posts.Posts = append(posts.Posts, *post)
	if err = posts.UpdatePost(postNum, false); err != nil {
		posts.Lg.Error("error edit post: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while edit post")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
}

// GetCreatePost shows clean form for new post.
func (c *MainController) GetCreatePost() {
	c.Data["BlogName"] = beego.AppConfig.String("BLOGNAME")
	c.TplName = "create.tpl"
}

// CreatePost create new post in DB.
func (c *MainController) CreatePost() {
	posts := models.NewPosts()

	post, err := c.decodePost()
	if err != nil {
		posts.Lg.Error("error while decoding new post body: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while decoding new post body")
		return
	}
	posts.Posts = append(posts.Posts, *post)
	if err = posts.CreatePost(); err != nil {
		posts.Lg.Error("error create new post: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while create new post")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
}

// decodePost is JSON decoder helper
func (c *MainController) decodePost() (*models.Post, error) {
	post := &models.Post{}
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(post); err != nil {
		return nil, err
	}
	return post, nil
}
