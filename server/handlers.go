package server

import (
	"blogMongo/models"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getTemplateHandler - возвращает шаблон
func (serv *Server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {

	blogs, err := models.GetAllBlogs(serv.ctx, serv.db)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	serv.Page.Title = "Мой блог"
	serv.Page.Data = blogs
	serv.Page.Command = "index"

	if err := serv.dictionary["BLOGS"].ExecuteTemplate(w, "base", serv.Page); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

// addGetBlogHandler - добавление блога
func (serv *Server) addGetBlogHandler(w http.ResponseWriter, r *http.Request) {

	serv.Page.Title = "Добавление блога"
	serv.Page.Data = models.Blog{}
	serv.Page.Command = "new"

	if err := serv.dictionary["BLOG"].ExecuteTemplate(w, "base", serv.Page); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

// addBlogHandler - добавляет блог
func (serv *Server) addBlogHandler(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	blog := models.Blog{}
	err = json.Unmarshal(data, &blog)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	if err := blog.AddBlog(serv.ctx, serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
	http.Redirect(w, r, "/", 301)
}

// deleteBlogHandler - удаляет блог
func (serv *Server) deleteBlogHandler(w http.ResponseWriter, r *http.Request) {

	blogID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	blog := models.Blog{}
	blog.ID = blogID

	if err := blog.Delete(serv.ctx, serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
	http.Redirect(w, r, "/", 301)
}

// editBlogHandler - редактирование блога
func (serv *Server) editBlogHandler(w http.ResponseWriter, r *http.Request) {

	blogID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	blog, err := models.GetBlog(serv.ctx, serv.db, blogID)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	serv.Page.Title = "Редактирование"
	serv.Page.Data = blog
	serv.Page.Command = "edit"

	if err := serv.dictionary["BLOG"].ExecuteTemplate(w, "base", serv.Page); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

// putBlogHandler - обновляет блог
func (serv *Server) putBlogHandler(w http.ResponseWriter, r *http.Request) {

	blogID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	blog := models.Blog{}
	err = json.Unmarshal(data, &blog)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	blog.ID = blogID

	if err := blog.UpdateBlog(serv.ctx, serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
	http.Redirect(w, r, "/", 301)
}
