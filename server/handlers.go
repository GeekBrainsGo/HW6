package server

import (
	"blogMongo/models"
	"net/http"
)

// getTemplateHandler - возвращает шаблон
func (serv *Server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {

	blogs, err := models.GetAllBlogItems(serv.ctx, serv.db)
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
