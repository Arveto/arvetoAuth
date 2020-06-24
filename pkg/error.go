// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"net/http"
	"strings"
)

type errorResponse struct {
	Code  int
	URL   string
	Error string
}

func (s *Server) Error(w http.ResponseWriter, r *http.Request, err string, code int) {
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.WriteHeader(code)
		s.errorPage.Execute(w, errorResponse{
			Code:  code,
			URL:   r.URL.String(),
			Error: err,
		})
	} else {
		http.Error(w, err, code)
	}
}
