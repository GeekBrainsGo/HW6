package server

import (
	"blogMongo/models"
	"net/http"

	"github.com/go-chi/chi"
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

// editBlogHandler - редактирование блога
func (serv *Server) editBlogHandler(w http.ResponseWriter, r *http.Request) {

	blogID := chi.URLParam(r, "id")

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
