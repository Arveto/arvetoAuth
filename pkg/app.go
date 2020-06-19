// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./public"
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

// Generate a JWT for a specific app.
func (s *Server) authUser(w http.ResponseWriter, r *http.Request) {
	var app application
	if s.db.GetS("app:"+r.URL.Query().Get("app"), &app) {
		http.Error(w, "Need app params in URL", http.StatusBadRequest)
		return
	}

	to := app.URL + "?"
	if r := r.URL.Query().Get("r"); r != "" {
		to += "r=" + r + "&"
	}

	if u := s.getUser(r); u != nil {
		if u.Level < public.LevelVisitor {
			http.Error(w, "Lowest level", http.StatusForbidden)
			return
		}
		jwt, err := u.ToJWT(s.key, app.ID)
		if err != nil {
			http.Error(w, "Generate JWT error: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, to+"jwt="+jwt, http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "No user (no login dev)", 400)
	}
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

func (s *Server) appAdd(w http.ResponseWriter, r *http.Request)  {}
func (s *Server) appRm(w http.ResponseWriter, r *http.Request)   {}
func (s *Server) appEdit(w http.ResponseWriter, r *http.Request) {}
