// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./public"
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/HuguesGuilleus/go-db.v1"
	"github.com/HuguesGuilleus/go-parsersa"
	"github.com/HuguesGuilleus/static.v2"
	"git.mills.io/prologic/bitcask"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"path/filepath"
)

// The option to create a server
type Option struct {
	Dev bool   // enable dev mode
	URL string // The URL of the server
	DB  string // path to the DB
	Key string // private key file
	// Mail options, need to be processed
	MailHost     string
	MailLogin    string
	MailPassword string
}

// One server. Use Option to create it.
type Server struct {
	db  *db.DB
	mux http.ServeMux
	key *rsa.PrivateKey
	url string
	// The template to send error response.
	errorPage *template.Template
	// The template to server /visit/ticket
	visitPage *template.Template
	// Mail options ready to use
	mailAuth  smtp.Auth
	mailHost  string
	mailLogin string
	nbAdmin   int
	// Default avatar
	avatarDefault http.HandlerFunc
}

func Create(opt Option) *Server {
	if opt.DB == "" {
		opt.DB = filepath.Join("data", "db")
	}

	k, err := parsersa.PrivFile(opt.Key)
	if err != nil {
		log.Fatal(err)
	}

	serv := &Server{
		url: opt.URL,
		db:  db.New(opt.DB, bitcask.WithMaxValueSize(3_000_000)),
		key: k,
		mailAuth: smtp.PlainAuth("",
			opt.MailLogin,
			opt.MailPassword,
			opt.MailHost),
		mailLogin: opt.MailLogin,
		mailHost:  opt.MailHost + ":smtp",
	}

	if opt.Dev {
		log.Println("[DEV MODE ENABLE]")
		static.Dev = true
		serv.loadDefaultUsers()
		serv.defaultApp()
		serv.mux.HandleFunc("/!users", serv.GodUsers)
		serv.mux.HandleFunc("/!login", serv.GodLogin)
	}

	serv.errorPage = static.TemplateHTML(nil, "front/error.html")
	serv.visitPage = static.TemplateHTML(nil, "front/visit.html")
	serv.avatarDefault = static.WebP(nil, "front/img/defaultUser.webp")

	go func() {
		serv.db.ForS("user:", 0, 0, nil, func(_ string, u *User) {
			if u.Level == public.LevelAdmin {
				serv.nbAdmin++
			}
		})
	}()

	// Static Handlers
	serv.mux.Handle("/style.css", static.Css(nil, "front/style.css"))
	serv.mux.Handle("/app.js", static.Js(nil, "front/app.js"))
	serv.mux.Handle("/favicon.ico", static.WebP(nil, "front/img/favicon.webp"))
	serv.mux.Handle("/arveto.webp", static.WebP(nil, "front/img/arveto.webp"))
	serv.mux.Handle("/alberto.webp", static.WebP(nil, "front/img/alberto.webp"))
	index := static.Html(nil, "front/index.html")
	pub := static.Html(nil, "front/public.html")
	serv.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			serv.Error(w, r, "Not Found", 404)
			return
		}

		if u := serv.getUser(r); u == nil {
			pub.ServeHTTP(w, r)
			return
		} else if u.Level < public.LevelVisitor {
			serv.Error(w, r, "Accréditation trop faible", http.StatusForbidden)
			return
		}
		index.ServeHTTP(w, r)
	})

	serv.handleLevel("/app/add", public.LevelAdmin, serv.appAdd)
	serv.handleLevel("/app/edit/name", public.LevelAdmin, serv.appEditName)
	serv.handleLevel("/app/edit/url", public.LevelAdmin, serv.appEditURL)
	serv.handleLevel("/app/list", public.LevelStd, serv.appList)
	serv.handleLevel("/app/rm", public.LevelAdmin, serv.appRm)
	serv.handleLevel("/auth", public.LevelCandidate, serv.authUser)
	serv.handleLevel("/avatar/edit", public.LevelStd, serv.avatarEdit)
	serv.handleLevel("/avatar/get", public.LevelCandidate, serv.avatarGet)
	serv.handleLevel("/avatar/invite", public.LevelCandidate,
		static.WebP(nil, "front/img/invite.webp"))
	serv.handleLevel("/publickey", public.LevelCandidate, static.File(serv.pubkey(), "", "text/plain", nil))
	serv.handleLevel("/log/count", public.LevelStd, serv.logCount)
	serv.handleLevel("/log/list", public.LevelStd, serv.logList)
	serv.handleLevel("/login/", public.LevelCandidate, static.Html(nil, "front/login.html"))
	serv.handleLevel("/login/from/github/", public.LevelCandidate, serv.loginFromGithub)
	serv.handleLevel("/login/with/github/", public.LevelCandidate, serv.loginWithGithub)
	serv.handleLevel("/login/from/google/", public.LevelCandidate, serv.loginFromGoogle)
	serv.handleLevel("/login/with/google/", public.LevelCandidate, serv.loginWithGoogle)
	serv.handleLevel("/logout", public.LevelCandidate, serv.logout)
	serv.handleLevel("/me", public.LevelCandidate, serv.getMe)
	serv.handleLevel("/sendmail", public.LevelAdmin, serv.testMail)
	serv.handleLevel("/user/edit/email", public.LevelStd, serv.userEditEmail)
	serv.handleLevel("/user/edit/level", public.LevelAdmin, serv.userEditLevel)
	serv.handleLevel("/user/edit/pseudo", public.LevelStd, serv.userEditPseudo)
	serv.handleLevel("/user/list", public.LevelCandidate, serv.userList)
	serv.handleLevel("/user/rm/me", public.LevelCandidate, serv.userRmMe)
	serv.handleLevel("/user/rm/other", public.LevelAdmin, serv.userRmOther)
	serv.handleLevel("/visit/add", public.LevelAdmin, serv.visitAdd)
	serv.handleLevel("/visit/list", public.LevelStd, serv.visitList)
	serv.handleLevel("/visit/ticket", public.LevelCandidate, serv.visitTicket)

	return serv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[REQ]", r.RequestURI)
	w.Header().Add("Server", "Arveto auth server")
	s.mux.ServeHTTP(w, r)
}

// Send JavaScrip to make redirection with cookie.
func redirection(w http.ResponseWriter, to string) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="utf-8"></head><body><a href="` + to + `">Redirect</a><script>document.location.replace('` + to + `');</script></body></html>`))
}

// Return the publix key used to sign JWT.
func (s *Server) pubkey() []byte {
	block := pem.Block{
		Type:  "BEGIN PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&s.key.PublicKey),
	}
	buff := bytes.Buffer{}
	pem.Encode(&buff, &block)
	return buff.Bytes()
}
