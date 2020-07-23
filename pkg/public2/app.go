// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package public

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"github.com/HuguesGuilleus/go-parsersa"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// The application information. Fill all fields before use it.
type App struct {
	// The address of the provider
	Provider string
	// The public key of the auth provider
	PublicKey *rsa.PublicKey
	// The audience field in JWT.
	Audience string
	// The name of the cookie who store the JWT
	// by default it's `auth`
	Cookie string
	//
	Mux http.ServeMux
	// A variant of http.Error used to send the error to the client. By default
	// its a binding of http.Error.
	Error func(w http.ResponseWriter, r *http.Request, err string, code int)
}

// Create a new app. id is the the id of the this application in the provider.
// provider is the url of this provider.
//
// defaultHandler registers the handler for login and logout
func NewApp(id, provider string, defaultHandler bool) (*App, error) {
	if !strings.HasSuffix(provider, "/") {
		provider += "/"
	}

	rep, err := http.Get(provider + "publickey")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}
	k, err := parsersa.Public(body)
	if err != nil {
		return nil, err
	}

	app := &App{
		Audience:  id,
		PublicKey: k,
		Cookie:    "auth",
		Provider:  provider,
		Error: func(w http.ResponseWriter, r *http.Request, err string, code int) {
			http.Error(w, err, code)
		},
	}

	if defaultHandler {
		app.Mux.HandleFunc("/login", app.Login(""))
		app.Mux.HandleFunc("/logout", app.Logout(""))
	}

	return app, nil
}

// Custom request. The user can be nil.
type Request struct {
	http.Request
	User *UserInfo
}

// HTTP.Handler with the custom request
type Handler interface {
	ServeHTTP(http.ResponseWriter, *Request)
}

// Add a handler to a.Mux. The level of the user must be over level.
//
// If the level is strict strict less than LevelCandidate, the user can be nil
// or with a lower level.
func (a *App) Handle(pattern string, level UserLevel, handler Handler) {
	a.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		u := a.User(r)
		if level > LevelCandidate {
			if u == nil {
				a.Error(w, r, "You are not logged", http.StatusUnauthorized)
				return
			} else if u.Level < level {
				a.Error(w, r, "Your level is too low; you need the level: "+level.String(), http.StatusForbidden)
				return
			}
		}
		handler.ServeHTTP(w, &Request{
			Request: *r,
			User:    u,
		})
	})
}

// HandleFunc with Request inplace of http.Request
type HandleFunc func(w http.ResponseWriter, r *Request)

func (f HandleFunc) ServeHTTP(w http.ResponseWriter, r *Request) { f(w, r) }

// Like App.Handle with a function inplace of a Handler.
func (a *App) HandleFunc(p string, l UserLevel, f func(w http.ResponseWriter, r *Request)) {
	a.Handle(p, l, HandleFunc(f))
}

// Binding of a.Mux.ServeHTTP
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Mux.ServeHTTP(w, r)
}

// Return the address top the provider to auth. r is the address to go after
// autentification work.
func (a *App) ProviderAuth(r string) string {
	out := a.Provider + "auth?app=" + a.Audience
	if r != "" {
		return out + "&r=" + tob64S(r)
	}
	return out
}

// Return a handler for the login. This handler takes the jwt from the URL and
// save it a cookie. Finaly it redirect the client to the params r or destination.
func (a *App) Login(destination string) http.HandlerFunc {
	if destination == "" {
		destination = "/"
	}

	model := http.Cookie{
		Name:     a.Cookie,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   24 * 60 * 60,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		j := r.URL.Query().Get("jwt")
		if _, err := a.FromJWT(j); err != nil {
			a.Error(w, r, "JWT error: "+err.Error(), http.StatusBadRequest)
			return
		}

		c := model
		c.Value = j
		w.Header().Add("Set-Cookie", c.String())

		to := fromb64(r.URL.Query().Get("r"))
		if to == "" {
			to = destination
		}
		redirection(w, to)
	}
}

// Remove the cookie and redirect the client to the variable to (or "/" if empty)
func (a *App) Logout(to string) http.HandlerFunc {
	if to == "" {
		to = "/"
	}

	c := (&http.Cookie{
		Name:     a.Cookie,
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}).String()

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", c)
		redirection(w, to)
	}
}

// Get the user from a request cookie. If error occure, the function
// return nil. If you want the error details, use App.FromJWT inplace.
func (a *App) User(r *http.Request) *UserInfo {
	c, err := r.Cookie(a.Cookie)
	if err != nil {
		return nil
	}

	u, err := a.FromJWT(c.Value)
	if err != nil {
		return nil
	}
	return u
}

var (
	JWTWrongSyntax     = errors.New("JWT wrong syntax")
	JWTWrongSyntaxHead = errors.New("JWT wrong syntax in head")
	JWTWrongHead       = errors.New("JWT wrong head")
	JWTWrongAudience   = errors.New("This JWT is made for an other audience")
	JWTOutDate         = errors.New("This JWT is out date")
	JWTEmpty           = errors.New("JWT is empty")
)

func (a *App) FromJWT(j string) (*UserInfo, error) {
	if j == "" {
		return nil, JWTEmpty
	}
	parts := strings.SplitN(j, ".", 3)
	if len(parts) != 3 {
		return nil, JWTWrongSyntax
	}

	// Check the head
	h := struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}{}
	json.Unmarshal([]byte(fromb64(parts[0])), &h)
	if h.Alg != "RS256" || h.Typ != "JWT" {
		return nil, JWTWrongHead
	}

	// Check the signature
	sig := []byte(fromb64(parts[2]))
	hash := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if err := rsa.VerifyPKCS1v15(a.PublicKey, crypto.SHA256, hash[:], sig); err != nil {
		return nil, err
	}

	// Get Body
	var body jwtBody
	if err := json.Unmarshal([]byte(fromb64(parts[1])), &body); err != nil {
		return nil, err
	}

	// Check Body
	if body.Audience != a.Audience {
		return nil, JWTWrongAudience
	}
	if body.ExpirationTime < time.Now().Unix() {
		return nil, JWTOutDate
	}

	return &body.UserInfo, nil
}
