package main

import (
	"HW6-master/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	STATIC_DIR = "./www/"
)

func main() {
	r := chi.NewRouter()
	lg := logrus.New()

	ctx := context.Background()
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	_ = client.Connect(ctx)
	db := client.Database("blog")

	serv := Server{
		lg:    lg,
		db:    db,
		ctx:   ctx,
		Title: "BLOG",
		Posts: models.Posts {
			{
				ID:		   1,
				Title:     "Пост 1",
				Text:      "Очень интересный текст",
				Labels:    []string{"путешестве", "отдых"},
			},
			{
				ID:		   2,
				Title:     "Пост 2",
				Text:      "Второй очень интересный текст",
				Labels:    []string{"домашка", "golang"},
			},
			{
				ID:		   3,
				Title:     "Пост 3",
				Text:      "Третий очень интересный текст",
				Labels:    []string{},
			},
		},
	}

	fileServer := http.FileServer(http.Dir(STATIC_DIR))
	r.Handle("/static/*", fileServer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.HandleGetIndex)
		r.Get("/post/{id}", serv.HandleGetPost)
		r.Get("/post/create", serv.HandleGetEditPost)
		r.Get("/post/{id}/edit", serv.HandleGetEditPost)
	})

	r.Route("/api/v1/", func(r chi.Router) {
		r.Post("/post/create", serv.HandleEditPost)
		r.Post("/post/{id}/edit", serv.HandleEditPost)
	})

	lg.Info("server is start")
	http.ListenAndServe(":8080", r)
}

type Server struct {
	lg    *logrus.Logger
	db    *mongo.Database
	ctx   context.Context
	Title string
	Posts models.Posts
}

func (serv *Server) AddOrUpdatePost(newPost models.Post) (models.Post) {

	updPost, err := newPost.Update(serv.ctx, serv.db)
	if err != nil {
		newPost, _ := newPost.Create(serv.ctx, serv.db)
		return *newPost
	}

	return *updPost
}


func (serv *Server) HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/templates/index.gohtml")
	data, _ := ioutil.ReadAll(file)

	posts, err := models.GetPosts(serv.ctx, serv.db)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetIndex")
		posts = models.Posts{}
	}

	serv.Posts = posts

	templ := template.Must(template.New("page").Parse(string(data)))
	err = templ.ExecuteTemplate(w, "page", serv)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetIndexTemplate")
	}
}

func (serv *Server) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/templates/post.gohtml")
	data, _ := ioutil.ReadAll(file)

	postIDStr := chi.URLParam(r, "id")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)

	post, err := models.GetPost(uint(postID), serv.ctx, serv.db)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetPost")
		post = &models.Post{}
	}

	fmt.Println(post.Text)

	templ := template.Must(template.New("page").Parse(string(data)))
	err = templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetPostTemplate")
	}
}

func (serv *Server) HandleGetEditPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/templates/edit_post.gohtml")
	data, _ := ioutil.ReadAll(file)

	postIDStr := chi.URLParam(r, "id")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)

	post, err := models.GetPost(uint(postID), serv.ctx, serv.db)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetEditPost")
		post = &models.Post{}
	}

	templ := template.Must(template.New("page").Parse(string(data)))
	err = templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetEditPostTemplate")
	}
}

func (serv *Server) HandleEditPost(w http.ResponseWriter, r *http.Request) {

	/*
	{"id":4, "Title":"Пост 4", "Text":"Новый очень интересный текст", "Labels":["l1","l2"]}
	*/

	decoder := json.NewDecoder(r.Body)
	var inPostItem models.Post
	err := decoder.Decode(&inPostItem)
	if err != nil {
		serv.lg.WithError(err).Error("HandleEditPost")
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	newPost := serv.AddOrUpdatePost(inPostItem)
	respondWithJSON(w, http.StatusOK, newPost)
}

// respondWithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

