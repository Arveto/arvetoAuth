// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./public"
	"encoding/json"
	"net/http"
)

// All the informations about an user
type User struct {
	public.UserInfo
	Cookie string
	// Password string
}

/* EDIT USER */

func (s *Server) userEditName(w http.ResponseWriter, r *http.Request) {
	s.usersEdit(w, r, func(u *User, v string) {
		u.Name = v
	})
}
func (s *Server) userEditEmail(w http.ResponseWriter, r *http.Request) {
	s.usersEdit(w, r, func(u *User, v string) {
		u.Email = v
	})
}
func (s *Server) usersEdit(w http.ResponseWriter, r *http.Request, edit func(*User, string)) {
	if r.Method != "PATCH" {
		http.Error(w, "Need a PATCH Method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
		http.Error(w, "Expected `Content-Type: text/plain; charset=utf-8`",
			http.StatusUnsupportedMediaType)
		return
	}

	data := make([]byte, 100, 100)
	if n, _ := r.Body.Read(data); n == 0 {
		http.Error(w, "Expected a body\r\n", http.StatusBadRequest)
		return
	} else {
		data = data[:n]
	}

	u := s.getUser(r)
	edit(u, string(data))
	s.db.SetS("user:"+u.Login, u)
}

/* GET USER AND USER LEVEL */

// Add a http.HandlerFunc to the server.mux. If a client with a lowest level
// or without authentification, the request are rejected.
func (s *Server) handleLevel(pattern string, l public.UserLevel, h http.HandlerFunc) {
	errLevel := "Required a highest level (" + l.String() + ")"
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if u := s.getUser(r); u == nil {
			http.Error(w, "Need authentification", http.StatusUnauthorized)
			return
		} else if u.Level < l {
			http.Error(w, errLevel, http.StatusForbidden)
			return
		}
		h(w, r)
	})
}

// Send public information about the current user
func (s *Server) getMe(w http.ResponseWriter, r *http.Request) {
	me := s.getUser(r)
	if me == nil {
		http.Error(w, "Who are you?", http.StatusUnauthorized)
		return
	}
	j, _ := json.Marshal(&me.UserInfo)
	w.Header().Add("Content-Type", "application/json")
	w.Write(j)
}

// Get the user from cookies and check-it with the DB.
func (s *Server) getUser(r *http.Request) *User {
	idCookie, _ := r.Cookie("id")
	creditCookie, _ := r.Cookie("credit")
	if idCookie == nil || creditCookie == nil {
		return nil
	}

	var u User
	if s.db.GetS("user:"+idCookie.Value, &u) {
		return nil
	}

	// TODO: check identity with cookie.

	return &u
}

func setCookie(w http.ResponseWriter, name, value string) {
	w.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}).String())
}

/* !!! DEVELOPMENT FUNCTIONS !!! */
// À utiliser QUE pour le dévelopement

func (s *Server) GodUsers(w http.ResponseWriter, r *http.Request) {
	list := make([]User, 0)
	s.db.ForS("user:", 0, 0, nil, func(_ string, u User) {
		list = append(list, u)
	})

	j, _ := json.Marshal(list)
	w.Header().Add("Content-Type", "application/json")
	w.Write(j)
}

func (s *Server) GodLogin(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("u")
	var u User
	if s.db.GetS("user:"+id, &u) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	u.Cookie = "pass"
	s.db.SetS("user:"+id, &u)
	setCookie(w, "id", r.URL.Query().Get("u"))
	setCookie(w, "credit", "pass")

	http.Redirect(w, r, "/me", http.StatusTemporaryRedirect)
}
