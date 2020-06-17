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

// Send public information about the current user
func (s *Server) getMe(w http.ResponseWriter, r *http.Request) {
	// if userOk(w, r) {
	// 	return
	// }
	//
	// u := User{}
	// if Users.GetS("user:"+r.URL.Query().Get("u"), &u) {
	// 	http.NotFound(w, r)
	// 	return
	// }

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
	return &u
}

/* !!! DEVELOPMENT FUNCTIONS !!! */
// À utiliser QUE pour le dévelopement

func (s *Server) GodUsers(w http.ResponseWriter, r *http.Request) {
	list := make([]User, 0)
	s.db.ForS("user:", func(_ string, u User) {
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

func setCookie(w http.ResponseWriter, name, value string) {
	w.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}).String())
}
