// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./public"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// All the informations about an user
type User struct {
	public.UserInfo
	Cookie map[string]time.Time
	// Password string
}

func (s *Server) userEditLevel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		s.Error(w, r, "Need a PATCH Method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		s.Error(w, r, "Expected `Content-Type: application/json`",
			http.StatusUnsupportedMediaType)
		return
	}

	var to public.UserLevel
	content, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err := to.UnmarshalJSON(content); err != nil {
		s.Error(w, r, "Parsing error: "+err.Error(), http.StatusBadRequest)
		return
	}

	var u User
	if s.db.GetS("user:"+r.URL.Query().Get("u"), &u) {
		s.Error(w, r, "User not found", http.StatusNotFound)
		return
	} else if u.Level == to {
		s.Error(w, r, "Same Level", http.StatusBadRequest)
		return
	}

	admin := s.getUser(r)
	if u.ID == admin.ID {
		s.Error(w, r, "You can't edit your own level", http.StatusBadRequest)
		return
	}

	s.logAdd(admin, "/user/edit/level", u.ID,
		u.Level.String()+" --> "+to.String())

	if u.Level == public.LevelAdmin {
		s.nbAdmin--
	} else if to == public.LevelAdmin {
		s.nbAdmin++
	}

	go s.sendMail(u.Email, "Changement d'accréditation",
		"Votre niveau d'accréditation à changé:\r\n"+
			u.Level.String()+" --> "+to.String()+"\r\n\r\n"+
			"Réalisé par: "+admin.Pseudo+
			"\r\n\r\n")

	u.Level = to
	s.db.SetS("user:"+u.ID, &u)
}

/* EDIT USER */

func (s *Server) userEditPseudo(w http.ResponseWriter, r *http.Request) {
	s.usersEdit(w, r, func(u *User, v string) {
		u.Pseudo = v
	})
}
func (s *Server) userEditEmail(w http.ResponseWriter, r *http.Request) {
	s.usersEdit(w, r, func(u *User, v string) {
		u.Email = v
	})
}
func (s *Server) usersEdit(w http.ResponseWriter, r *http.Request, edit func(*User, string)) {
	if r.Method != "PATCH" {
		s.Error(w, r, "Need a PATCH Method", http.StatusBadRequest)
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

	u := s.getUser(r)
	edit(u, string(data))
	s.db.SetS("user:"+u.ID, u)
}

// The user remove itself.
func (s *Server) userRmMe(w http.ResponseWriter, r *http.Request) {
	u := s.getUser(r)
	if u == nil {
		s.Error(w, r, "Vous n'êtes pas identifié", http.StatusUnauthorized)
		return
	}

	s.logAdd(u, "/user/rm/me")
	s.db.DeleteS("user:" + u.ID)

	s.Error(w, r, "Supprétion réussi", http.StatusOK)
}

// An administrator remove an other user.
func (s *Server) userRmOther(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("u")
	if id == "" {
		s.Error(w, r, "Need an user (?u=userID) in params", http.StatusBadRequest)
		return
	}
	id = "user:" + id

	var u User
	if s.db.GetS(id, &u) {
		s.Error(w, r, "User not Found", http.StatusNotFound)
		return
	}

	s.logAdd(s.getUser(r), "/user/rm/other", u.ID)
	s.db.DeleteS(id)
	s.Error(w, r, "User deleted", http.StatusOK)
}

/* GET USER AND USER LEVEL */

// Add a http.HandlerFunc to the server.mux. If a client with a lowest level
// or without authentification, the request are rejected.
func (s *Server) handleLevel(pattern string, l public.UserLevel, h http.HandlerFunc) {
	if l == public.LevelCandidate {
		s.mux.HandleFunc(pattern, h)
		return
	}

	errLevel := "Required a highest level (" + l.String() + ")"
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if u := s.getUser(r); u == nil {
			s.Error(w, r, "Need authentification", http.StatusUnauthorized)
			return
		} else if u.Level < l {
			s.Error(w, r, errLevel, http.StatusForbidden)
			return
		}
		h(w, r)
	})
}

// List the users. It depend of the user who make the requet.
func (s *Server) userList(w http.ResponseWriter, r *http.Request) {
	u := s.getUser(r)
	if u == nil {
		s.Error(w, r, "Need authentification", http.StatusUnauthorized)
		return
	}

	var filter func(*User) bool
	switch u.Level {
	case public.LevelStd:
		filter = func(u *User) bool { return u.Level >= public.LevelStd }
	case public.LevelAdmin:
		filter = func(*User) bool { return true }
	default:
		s.Error(w, r, "Required a stantard level", http.StatusForbidden)
		return
	}

	all := make([]public.UserInfo, 0)
	s.db.ForS("user:", 0, 0, nil, func(_ string, u *User) {
		if filter(u) {
			all = append(all, u.UserInfo)
		}
	})

	w.Header().Add("Content-Type", "application/json")
	j, _ := json.Marshal(all)
	w.Write(j)
}

// Send public information about the current user
func (s *Server) getMe(w http.ResponseWriter, r *http.Request) {
	me := s.getUser(r)
	if me == nil {
		s.Error(w, r, "Who are you?", http.StatusUnauthorized)
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

	// Remove old Cookie
	save := false
	for k, t := range u.Cookie {
		if t.Before(time.Now()) {
			delete(u.Cookie, k)
			save = true
		}
	}
	if save {
		s.db.SetS("user:"+idCookie.Value, &u)
	}

	// The if the actual cookie is set.
	if u.Cookie[creditCookie.Value].IsZero() {
		return nil
	}

	return &u
}

// Create a new cookie for a user
func (s *Server) setCookie(w http.ResponseWriter, r *http.Request, u *User) {
	v := make([]byte, 15, 15)
	rand.Read(v)
	c := base64.RawStdEncoding.EncodeToString(v)
	if u.Cookie == nil {
		u.Cookie = make(map[string]time.Time, 1)
	}
	u.Cookie[c] = time.Now().Add(time.Hour * time.Duration(6))
	s.db.SetS("user:"+u.ID, u)

	w.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     "credit",
		Value:    c,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}).String())
	w.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     "id",
		Value:    u.ID,
		Path:     "/",
		HttpOnly: true,
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
		s.Error(w, r, "User not found", http.StatusNotFound)
		return
	}

	s.setCookie(w, r, &u)
	http.Redirect(w, r, "/me", http.StatusFound)
}
