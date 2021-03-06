// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/json"
	"net/http"
)

type application struct {
	ID   string
	Name string
	URL  string
}

// Création de valeurs par défaut
func (s *Server) defaultApp() {
	s.db.SetS("app:ex", application{
		ID:   "ex",
		Name: "Example for the dev",
		URL:  "http://localhost:9000/login",
	})
}

// List the applications.
func (s *Server) appList(w http.ResponseWriter, r *http.Request) {
	all := make([]application, 0)
	s.db.ForS("app:", 0, 0, nil, func(_ string, a application) {
		all = append(all, a)
	})
	j, _ := json.Marshal(all)
	w.Header().Add("Content-Type", "application/json")
	w.Write(j)
}

// Create a new application.
func (s *Server) appAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.Error(w, r, "Need a POST Method", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	k := "app:" + id
	if id == "" {
		s.Error(w, r, "Need and `id` in params\r\n", http.StatusBadRequest)
		return
	} else if !s.db.UnknownS(k) {
		s.Error(w, r, "This app already exist\r\n", http.StatusConflict)
		return
	}

	s.db.SetS(k, application{ID: id})
	s.logAdd(s.getUser(r), "/app/add", id)
}

// Remove an application
func (s *Server) appRm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.Error(w, r, "Need a POST Method", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	k := "app:" + id
	if id == "" {
		s.Error(w, r, "Need and `id` in params\r\n", http.StatusBadRequest)
		return
	} else if s.db.UnknownS(k) {
		s.Error(w, r, "This app does not exist\r\n", http.StatusNotFound)
		return
	}

	s.db.DeleteS(k)
	s.logAdd(s.getUser(r), "/app/rm", id)
}

// Change the URL of an application.
func (s *Server) appEditURL(w http.ResponseWriter, r *http.Request) {
	s.appEdit(w, r, func(app *application, v string) {
		app.URL = v
		s.logAdd(s.getUser(r), "/app/edit/url", app.ID, v)
	})
}

// Change the name of the application.
func (s *Server) appEditName(w http.ResponseWriter, r *http.Request) {
	s.appEdit(w, r, func(app *application, v string) {
		app.Name = v
		s.logAdd(s.getUser(r), "/app/edit/name", app.ID, v)
	})
}

// appEdit get the app and the data from the body and call the function edit
// to edit the application.
func (s *Server) appEdit(w http.ResponseWriter, r *http.Request, edit func(*application, string)) {
	if r.Method != "PATCH" {
		s.Error(w, r, "Need a PATCH Method", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	k := "app:" + id
	app := application{}
	if id == "" {
		s.Error(w, r, "Need and `id` in params\r\n", http.StatusBadRequest)
		return
	} else if s.db.GetS(k, &app) {
		s.Error(w, r, "This app does not exist\r\n", http.StatusNotFound)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
		s.Error(w, r, "Expected `Content-Type: text/plain; charset=utf-8`",
			http.StatusUnsupportedMediaType)
		return
	}

	data := make([]byte, 100, 100)
	if n, _ := r.Body.Read(data); n == 0 {
		s.Error(w, r, "Expected a body\r\n", http.StatusBadRequest)
		return
	} else {
		data = data[:n]
	}
	edit(&app, string(data))

	s.db.SetS(k, &app)
}
