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

func (s *Server) loginIn(w http.ResponseWriter, r *http.Request, app string) {
	http.Redirect(w, r, github.URL(s.url+"login/github/"+app+"/"+r.URL.Query().Get("r")),
		http.StatusTemporaryRedirect)
}

func (s *Server) loginGithub(w http.ResponseWriter, r *http.Request) {
	info, err := github.NewInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := "github:" + info.Login
	var u User
	if s.db.GetS("user:"+id, &u) {
		u.Login = id
		u.Name = info.Pseudo
		u.Email = info.Email
		u.Level = public.LevelCandidate
		s.db.SetS("user:"+id, &u)

		// TODO: get Avatar

		http.Error(w, "Vous êtes inscris! mais vous devez être accrédité pour continuer", http.StatusForbidden)
		return
	}

	s.setCookie(w, r, &u)
	http.Redirect(w, r,
		loginApp.ReplaceAllString(r.URL.Path, "/auth?app=$1&r=$2"),
		http.StatusFound)
}
