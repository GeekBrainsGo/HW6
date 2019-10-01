package server

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"serv/models"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type Server struct {
	lg            *logrus.Logger
	ctx           context.Context
	db            *mongo.Database
	rootDir       string
	templatesDir  string
	indexTemplate string
	Title         string
	Posts         models.PostItemSlice
}

// New - создаёт новый экземпляр сервера
func New(lg *logrus.Logger, rootDir string, ctx context.Context, db *mongo.Database) *Server {
	return &Server{
		lg:            lg,
		db:            db,
		rootDir:       rootDir,
		templatesDir:  "/templates",
		indexTemplate: "index.html",
		Posts:         models.PostItemSlice{},
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
		r.Get("/", serv.HandleGetIndexHtml)
		r.Get("/view/{postID}", serv.HandleGetPostHtml)
		r.Get("/edit/{postID}", serv.HandleGetEditHtml)
		r.Get("/new", serv.HandleGetNewHtml)
		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/posts/new", serv.postNewHandler)
			r.Delete("/posts/{id}", serv.deletePostHandler)
			r.Put("/posts/{id}", serv.putPostHandler)
		})
	})

}
