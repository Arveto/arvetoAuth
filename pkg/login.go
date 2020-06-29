// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./public"
	"./public/github"
	"./public/google"
	"net/http"
)

const newUser = "Vous êtes inscris! Mais vous devez être accrédité pour continuer."

// Redirect to an Github authentification service.
func (s *Server) loginWithGoogle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, google.URL("/"), http.StatusTemporaryRedirect)
}

// Redirect to an Google authentification service.
func (s *Server) loginWithGithub(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, github.URL(""), http.StatusTemporaryRedirect)
}

func (s *Server) loginFromGoogle(w http.ResponseWriter, r *http.Request) {
	s.loginFrom(w, r, google.User)
}

// Manage user from GitHub authentification service.
func (s *Server) loginFromGithub(w http.ResponseWriter, r *http.Request) {
	s.loginFrom(w, r, github.User)
}

// Manage user login from a provider.
func (s *Server) loginFrom(w http.ResponseWriter, r *http.Request, p public.Provider) {
	info, err := p(r)
	if err != nil {
		s.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	var u User
	if s.db.GetS("user:"+info.ID, &u) {
		defer s.db.SetS("user:"+info.ID, &u)
		u.UserInfo = *info

		go s.avatarFromURL(&u, info.Avatar)
		u.Avatar = s.url + "avatar/get?u=" + u.ID

		if s.nbAdmin > 0 {
			u.Level = public.LevelCandidate
			s.Error(w, r, newUser, http.StatusForbidden)
			return
		}
		u.Level = public.LevelAdmin
	} else if u.Level < public.LevelStd {
		s.Error(w, r, newUser, http.StatusForbidden)
		return
	}

	s.setCookie(w, r, &u)
	if authApp(r) != "" {
		redirection(w, "/auth")
	} else {
		redirection(w, "/")
	}
}

// Remove user's cookie.
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
