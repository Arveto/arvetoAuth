// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./db"
	"log"
	"net/http"
	"path/filepath"
)

// The option to create a server
type Option struct {
	DB string // path to the DB
}

func Create(opt Option) *Server {
	if opt.DB == "" {
		opt.DB = filepath.Join("data", "db")
	}
	serv := &Server{
		db: db.New(opt.DB),
	}

	// Remove this lines for production
	// serv.loadDefaultUsers()
	serv.mux.HandleFunc("/!users", serv.GodUsers)
	serv.mux.HandleFunc("/!login", serv.GodLogin)

	serv.mux.HandleFunc("/me", serv.getMe)

	return serv
}

// One server. Use Option to create it.
type Server struct {
	db  *db.DB
	mux http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[REQ]", r.URL.Path)
	w.Header().Add("Server", "Arveto auth server")
	s.mux.ServeHTTP(w, r)
}
