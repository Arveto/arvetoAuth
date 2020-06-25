// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./github"
	"./public"
	"net/http"
	"regexp"
)

var loginApp = regexp.MustCompile(`/login/\w+/(\w+)/(.*)`)

// Redirect to an authentification service.
func loginInHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, github.URL(""), http.StatusTemporaryRedirect)
}

// Redirect to a specific service and prepare the redirection to
// an specific application.
func (s *Server) loginIn(w http.ResponseWriter, r *http.Request, app string) {
	http.Redirect(w, r,
		github.URL(s.url+"login/github/"+app+"/"+r.URL.Query().Get("r")),
		http.StatusTemporaryRedirect)
}

// Manage user from GitHub authentification service.
func (s *Server) loginGithub(w http.ResponseWriter, r *http.Request) {
	info, err := github.NewInfo(r)
	if err != nil {
		s.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	id := "github:" + info.Login
	var u User
	if s.db.GetS("user:"+id, &u) {
		defer s.db.SetS("user:"+id, &u)
		u.ID = id
		u.Pseudo = info.Pseudo
		u.Email = info.Email

		u.Avatar = s.url + "avatar/get?u=" + u.ID
		go s.avatarFromURL(&u, info.Icon)

		if s.nbAdmin > 0 {
			u.Level = public.LevelCandidate
			s.Error(w, r, "Vous êtes inscris! mais vous devez être accrédité pour continuer", http.StatusForbidden)
			return
		}
		u.Level = public.LevelAdmin
	}

	s.setCookie(w, r, &u)

	if loginApp.MatchString(r.URL.Path) {
		redirection(w, loginApp.ReplaceAllString(r.URL.Path, "/auth?app=$1&r=$2"))
	} else {
		redirection(w, "/")
	}
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	u := s.getUser(r)
	if u == nil {
		s.Error(w, r, "Vous n'êtes pas identifié.", http.StatusBadRequest)
		return
	}
	c, _ := r.Cookie("credit")
	delete(u.Cookie, c.Value)
	s.db.SetS("user:"+u.ID, u)

	redirection(w, "/")
}
