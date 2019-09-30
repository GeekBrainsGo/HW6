package server

import (
	"blogMongo/models"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"path"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TEMPLATESDIR = "/templates"
)

// Server - объект сервера
type Server struct {
	ctx        context.Context
	lg         *logrus.Logger
	db         *mongo.Database
	rootDir    string
	Page       models.Page
	dictionary map[string]*template.Template
}

// New - создаёт новый экземпляр сервера
// func New(lg *logrus.Logger, rootDir string, db *mongo.Database) *Server {
func New(cont context.Context, lg *logrus.Logger, rootDir string, db *mongo.Database) *Server {

	// initial
	fileBase := path.Join(rootDir, TEMPLATESDIR, "base.html")
	fileMenu := path.Join(rootDir, TEMPLATESDIR, "menu.html")
	tMap := make(map[string]*template.Template)

	fileAdd := path.Join(rootDir, TEMPLATESDIR, "blogs.html")
	temp := template.Must(template.ParseFiles(fileBase, fileMenu, fileAdd))
	tMap["BLOGS"] = temp

	fileAdd = path.Join(rootDir, TEMPLATESDIR, "blog.html")
	temp = template.Must(template.ParseFiles(fileBase, fileMenu, fileAdd))
	tMap["BLOG"] = temp

	return &Server{
		ctx:        cont,
		lg:         lg,
		db:         db,
		rootDir:    rootDir,
		Page:       models.Page{},
		dictionary: tMap,
	}
}

// Start - запускает сервер
func (serv *Server) Start(addr string) error {
	r := chi.NewRouter()
	serv.bindRoutes(r)
	serv.lg.Debug("server is started ...")
	return http.ListenAndServe(addr, r)
}

func (serv *Server) bindRoutes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.getTemplateHandler)
		// r.Route("/api/v1", func(r chi.Router) {
		// 	r.Put("/blog/{id}", serv.putBlogHandler)
		// 	r.Get("/blog/{id}", serv.editBlogHandler)
		// 	r.Post("/blog", serv.addBlogHandler)
		// 	r.Get("/blog", serv.addGetBlogHandler)
		// 	r.Delete("/blog/{id}", serv.deleteBlogHandler)
		// })
	})
}

// SendErr - отправляет ошибку пользователю и логирует её
func (serv *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	serv.lg.WithField("data", obj).WithError(err).Error("server error")
	w.WriteHeader(code)
	errModel := models.ErrorModel{
		Code:     code,
		Err:      err.Error(),
		Desc:     "server error",
		Internal: obj,
	}
	data, _ := json.Marshal(errModel)
	w.Write(data)
}

// SendInternalErr - отправляет 500 ошибку
func (serv *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	serv.SendErr(w, err, http.StatusInternalServerError, obj)
}
