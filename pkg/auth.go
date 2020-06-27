// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./public"
	"net/http"
)

// Generate a JWT for a specific app.
func (s *Server) authUser(w http.ResponseWriter, r *http.Request) {
	var app application
	if a := authApp(r); a == "" || s.db.GetS("app:"+a, &app) {
		s.Error(w, r, "Need app params in URL or cookie 'app'",
			http.StatusBadRequest)
		return
	}

	redirect := r.URL.Query().Get("r")

	if u := s.getUser(r); u != nil {
		if u.Level < public.LevelStd {
			s.Error(w, r, "Lowest level", http.StatusForbidden)
			return
		}

		jwt, err := u.ToJWT(s.key, app.ID)
		if err != nil {
			s.Error(w, r, "Generate JWT error: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		to := app.URL + "?jwt=" + jwt
		if redirect != "" {
			to += "&r=" + redirect
		}
		authRmAuth(w)
		redirection(w, to)
	} else {
		authLogin(w, app.ID, redirect)
	}
}

// Get the app from app URL param or cookie.
func authApp(r *http.Request) string {
	if a := r.URL.Query().Get("app"); a != "" {
		return a
	}
	if c, _ := r.Cookie("app"); c != nil {
		return c.Value
	}
	return ""
}

// Get the redirect URL
func authRedirect(r *http.Request) string {
	if redirect := r.URL.Query().Get("r"); redirect != "" {
		return redirect
	}
	if c, _ := r.Cookie("redirect"); c != nil {
		return c.Value
	}
	return ""
}

// Redirect the user to the login page.
func authLogin(w http.ResponseWriter, app, r string) {
	w.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     "app",
		Value:    app,
		Path:     "/",
		HttpOnly: true,
	}).String())
	w.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     "redirect",
		Value:    r,
		Path:     "/",
		HttpOnly: true,
	}).String())
	redirection(w, "/login/")
}

// Supprime les cookie
func authRmAuth(w http.ResponseWriter) {
	w.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     "app",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}).String())
	w.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     "redirect",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}).String())
}
