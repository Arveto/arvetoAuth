// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./avatarreencode"
	"net/http"
)

func (s *Server) avatarGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("u")
	if id == "" {
		s.Error(w, r, "Give user with u in URL params", http.StatusBadRequest)
		return
	}

	img := s.db.GetSRaw("avatar:" + id)
	if len(img) == 0 {
		s.avatarDefault.ServeHTTP(w, r)
		return
	}

	w.Header().Add("Content-Type", "image/webp")
	w.Write(img)
}

func (s *Server) avatarEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		s.Error(w, r, "Use PATH method to set an avatar",
			http.StatusMethodNotAllowed)
		return
	}

	img, err := avatarreencode.Reencode(r.Body, r.Header.Get("Content-Type"))
	switch err {
	case nil:
	case avatarreencode.ImageTypeUnknown:
		s.Error(w, r, err.Error(), http.StatusUnsupportedMediaType)
		return
	default:
		s.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	s.db.SetSRaw("avatar:"+s.getUser(r).ID, img)
}

// Get the avatar from an URL. The error should no logged.
func (s *Server) avatarFromURL(u *User, from string) {
	if from == "" {
		return
	}

	rep, err := http.Get(from)
	if err != nil {
		return
	}
	img, err := avatarreencode.Reencode(rep.Body, rep.Header.Get("Content-Type"))
	if err != nil {
		return
	}
	s.db.SetSRaw("avatar:"+u.ID, img)
}
