package server

import (
	"encoding/json"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"serv/models"
	"text/template"
)

// HandleGetIndexHtml - возвращает главную страницу - index.html
func (serv *Server) HandleGetIndexHtml(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open(serv.rootDir + "/static/index.html")
	data, _ := ioutil.ReadAll(file)

	posts, errAllPosts := models.GetAllPosts(serv.ctx, serv.db)
	if errAllPosts != nil {
		serv.lg.WithError(errAllPosts).Error("GetAllPosts")
	}
	serv.Posts = posts

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", serv)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

// HandleGetPostHtml - возвращает страницу конкретного поста - post.html
func (serv *Server) HandleGetPostHtml(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")

	postFindId := -1
	for i := 0; i < len(serv.Posts); i++ {
		if serv.Posts[i].Idn == postIDStr {
			postFindId = i
		}
	}
	post := serv.Posts[postFindId]

	file, _ := os.Open(serv.rootDir + "/static/post.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

// HandleGetEditHtml - возвращает страницу редактирования поста - edit.html
func (serv *Server) HandleGetEditHtml(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")

	postFindId := 0
	for i := 0; i < len(serv.Posts); i++ {
		if serv.Posts[i].Idn == postIDStr {
			postFindId = i
		}
	}
	post := serv.Posts[postFindId]

	file, _ := os.Open(serv.rootDir + "/static/edit.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

// HandleGetNewHtml - возвращает страницу создания нового поста - new.html
func (serv *Server) HandleGetNewHtml(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open(serv.rootDir + "/static/new.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", serv)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

// postNewHandler - добавление нового поста
func (serv *Server) postNewHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)

	post := models.PostItem{}
	_ = json.Unmarshal(data, &post)
	post.Idn = uuid.NewV4().String()

	p, err := post.Insert(serv.ctx, serv.db)
	if err != nil {
		serv.lg.WithError(err).Error("DB insert")
		return
	}

	data, _ = json.Marshal(p)
	w.Write(data)
}

// deletePostHandler - удаляем пост
func (serv *Server) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	post := models.PostItem{Idn: postID}

	if _, err := post.Delete(serv.ctx, serv.db); err != nil {
		serv.lg.WithError(err).Error("DB delete")
		return
	}
}

// putPostHandler - обновляем пост
func (serv *Server) putPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	//fmt.Println(postID)
	data, _ := ioutil.ReadAll(r.Body)

	post := models.PostItem{}
	_ = json.Unmarshal(data, &post)
	post.Idn = postID

	data, _ = json.Marshal(post)
	w.Write(data)

	if _, err := post.Update(serv.ctx, serv.db); err != nil {
		serv.lg.WithError(err).Error("DB update")
		return
	}

}
