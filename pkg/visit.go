// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./public"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// One invitation ticket
type Visit struct {
	public.UserInfo
	Creation time.Time `json:"creation"`
	Author   string    `json:"author"` // Author ID
	App      string    `json:"app"`    // the app ID
	URL      string    `json:"url"`    // the URL to the application auth
}

// Create a new invitation.
func (s *Server) visitAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.Error(w, r, "Need a POST Method", http.StatusBadRequest)
		return
	} else if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		s.Error(w, r, "Expected `Content-Type: application/x-www-form-urlencoded`",
			http.StatusUnsupportedMediaType)
		return
	} else if err := r.ParseForm(); err != nil {
		s.Error(w, r, "Parse body error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	admin := s.getUser(r)
	v := Visit{
		UserInfo: public.UserInfo{
			ID:     time.Now().Format("2006-01-02:05.00000"),
			Pseudo: r.FormValue("pseudo"),
			Email:  r.FormValue("email"),
			Avatar: s.url + "avatar/invite",
			Level:  public.LevelVisitor,
		},
		Creation: time.Now().Truncate(time.Second),
		Author:   admin.ID,
		App:      r.FormValue("app"),
	}

	if v.Pseudo == "" || v.Email == "" || v.App == "" {
		s.Error(w, r, "Need: a pseudo, an email and a app.", http.StatusBadRequest)
		return
	}

	var app application
	if s.db.GetS("app:"+v.App, &app) {
		s.Error(w, r, "Not found application", http.StatusNotFound)
		return
	}

	jwt, err := v.UserInfo.ToJWT(s.key, app.ID)
	if err != nil {
		s.Error(w, r, "Generate JWT error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	v.URL = app.URL + "?jwt=" + jwt

	s.db.SetS("visit:"+v.ID, &v)
	s.logAdd(admin, "/visit/add", v.ID, v.App, v.Email)
	go s.sendMail(v.Email, "Ticket de visite", "Bonjour, Vous pouvez accéder au site\u00A0: "+app.Name+" pendant 25 heures via ce lien\u00A0:\r\n\r\n"+v.URL+"\r\n\r\nBonne journée,\r\n"+admin.Pseudo+".\r\n")

	w.Write([]byte("ok"))
}

const day time.Duration = 24 * time.Hour

// Send the invitation ticket
func (s *Server) visitTicket(w http.ResponseWriter, r *http.Request) {
	var v Visit
	if s.db.GetS("visit:"+r.URL.Query().Get("v"), &v) {
		s.Error(w, r, "Visit not found", http.StatusNotFound)
		return
	}

	url := v.URL
	if v.Creation.Add(day).Before(time.Now()) {
		url = ""
	}

	s.visitPage.Execute(w, struct{ Name, Email, URL string }{
		Name:  v.Pseudo,
		Email: v.Email,
		URL:   url,
	})
}

// List the invitation
func (s *Server) visitList(w http.ResponseWriter, r *http.Request) {
	prefix, err := getDuration(w, r, "visit:")
	if err {
		return
	}

	log.Printf("prefix: %#+v\n", prefix)

	list := make([]Visit, 0)
	s.db.ForS(prefix, 0, 0, nil, func(_ string, v Visit) {
		list = append(list, v)
	})

	j, _ := json.Marshal(list)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
