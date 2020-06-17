package auth

import (
	"./public"
	"encoding/base64"
	"net/http"
	"strings"
)

type application struct {
	ID   string
	Name string
	URL  string
}

// Création de valeurs par défaut
func (s *Server) defaultApp() {
	s.db.SetS("app:ex", application{
		ID:   "ex",
		Name: "Example for the dev",
		URL:  "http://localhost:9000/",
	})
}

func (s *Server) authUser(w http.ResponseWriter, r *http.Request) {
	var app application
	if s.db.GetS("app:"+r.URL.Query().Get("app"), &app) {
		http.Error(w, "Need app params in URL", http.StatusBadRequest)
		return
	}

	toBytes, _ := base64.RawURLEncoding.DecodeString(r.URL.Query().Get("r"))
	to := string(toBytes)
	if to == "" {
		to = app.URL
	} else if !strings.HasPrefix(app.URL, to) {
		http.Error(w, "Redirect params and app.URL have not the same begin.",
			http.StatusBadRequest)
		return
	}

	if u := s.getUser(r); u != nil {
		if u.Level < public.LevelVisitor {
			http.Error(w, "Lowest level", http.StatusForbidden)
			return
		}
		jwt, err := u.ToJWT(s.key, app.ID)
		if err != nil {
			http.Error(w, "Generate JWT error: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, to+jwt,
			http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "No user (no login dev)", 400)
	}
}
