// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
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
