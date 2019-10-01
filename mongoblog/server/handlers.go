package server

import (
	"encoding/json"
	"HW6/mongoblog/models"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi"
)

// templateHandle returns template.
func (s *Server) templateHandle(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "template")

	if templateName == "" {
		templateName = s.indexTemplate
	}

	file, err := os.Open(path.Join(s.rootDir, s.templatesDir, templateName))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.SendInternalErr(w, err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("").Parse(string(data))
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}

	posts, err := models.AllPosts(s.db)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}

	s.Page.Posts = posts

	if err := templ.Execute(w, s.Page); err != nil {
		s.SendInternalErr(w, err)
		return
	}
}

// postHandle adds new post.
func (s *Server) postHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()

	var err error
	post := &models.Post{}
	if err = json.NewDecoder(r.Body).Decode(post); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	err = post.Insert(s.db)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}
	json.NewEncoder(w).Encode(post)
}

// deleteHandle deletes a post.
func (s *Server) deleteHandle(w http.ResponseWriter, r *http.Request) {
	hex := chi.URLParam(r, "id")
	post := models.Post{Hex: hex}
	if _, err := post.Delete(s.db); err != nil {
		s.SendInternalErr(w, err)
		return
	}
}

// putHandle renew post.
func (s *Server) putHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	
	var err error
	post := &models.Post{}
	if err = json.NewDecoder(r.Body).Decode(post); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	if _, err = post.Update(s.db); err != nil {
		s.SendInternalErr(w, err)
		return
	}
}
