package webserver

import (
	"blog/app/models"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// WebServer ...
type WebServer struct {
	router   *chi.Mux
	logger   *logrus.Logger
	database *mgo.Session
}

func newServer(db *mgo.Session) *WebServer {
	serv := &WebServer{
		router:   chi.NewRouter(),
		logger:   logrus.New(),
		database: db,
	}

	serv.configureRouter()

	return serv
}

// Start ...
func Start(config *Config) error {
	db, err := newSession(config.DatabaseConnectionString)
	if err != nil {
		return err
	}

	defer db.Close()
	serv := newServer(db)
	return http.ListenAndServe(config.BindAddr, serv)
}

func newSession(dsnURL string) (*mgo.Session, error) {

	session, err := mgo.Dial(dsnURL)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (serv *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serv.router.ServeHTTP(w, r)
}

func (serv *WebServer) configureRouter() {
	//routes
	serv.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	serv.router.HandleFunc("/list", serv.postListHandle())

	serv.router.HandleFunc("/view/{postID}", serv.postViewHandle())

	serv.router.HandleFunc("/delete/{postID}", serv.postDeleteHandle())

	serv.router.HandleFunc("/create", serv.postCreateHandle())

}

func (serv *WebServer) postListHandle() http.HandlerFunc {

	type PageModel struct {
		Title string
		Data  interface{}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		conn := serv.database.DB("blog").C("posts")

		var posts models.PostItemsSlice

		err := conn.Find(bson.M{}).All(&posts)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		pageData := PageModel{
			Title: "Posts List",
			Data:  posts,
		}

		templ := template.Must(template.New("page").ParseFiles("./templates/blog/List.tpl", "./templates/common.tpl"))
		err = templ.ExecuteTemplate(w, "page", pageData)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

	}
}

func (serv *WebServer) postViewHandle() http.HandlerFunc {

	type PageModel struct {
		Title string
		Data  interface{}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		postID := chi.URLParam(r, "postID")

		conn := serv.database.DB("blog").C("posts")

		var post models.Post

		err := conn.Find(bson.M{"id": postID}).One(&post)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		pageData := PageModel{
			Title: "View Post",
			Data:  post,
		}

		templ := template.Must(template.New("page").ParseFiles("./templates/blog/View.tpl", "./templates/common.tpl"))
		err = templ.ExecuteTemplate(w, "page", pageData)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

	}
}

func (serv *WebServer) postDeleteHandle() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		postID := chi.URLParam(r, "postID")

		conn := serv.database.DB("blog").C("posts")

		err := conn.Remove(bson.M{"id": postID})
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return

	}
}

func (serv *WebServer) postCreateHandle() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var newPost models.Post

		newPost.ID = uuid.NewV4().String()
		newPost.Title = "New Post Title"
		newPost.Short = "Short body"
		newPost.Body = "Content body"

		conn := serv.database.DB("blog").C("posts")

		err := conn.Insert(newPost)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return

	}
}

func (serv *WebServer) errorAPI(w http.ResponseWriter, r *http.Request, code int, err error) {
	serv.respondJSON(w, r, code, map[string]string{"error": err.Error()})
}

func (serv *WebServer) respondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (serv *WebServer) respondWhithTemplate(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
