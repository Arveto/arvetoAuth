package public

import (
	"crypto/rsa"
	"encoding/base64"
	"github.com/HuguesGuilleus/go-parsersa"
	"io/ioutil"
	"net/http"
)

// The application information. Fill all fields before use it.
type App struct {
	// The name of the app
	ID string
	// The public key of the auth server
	PublicKey *rsa.PublicKey
	// The name of the cookie who store teh JWT
	// by default it's `auth`
	Cookie string
	// The address of the provider
	Server string
}

// Create a new app. id is the the id of the this application in the provider.
// serv is the url of this provider.
func NewApp(id, server string) (*App, error) {
	if server[len(server)-1] != '/' {
		server += "/"
	}

	rep, err := http.Get(server + "/publickey")
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

	return &App{
		ID:        id,
		PublicKey: k,
		Cookie:    "auth",
		Server:    server,
	}, nil
}

// Get the user from a request cookie. If error occure, the function
// return nil. If you want the error details, use FromJWT inplace.
func (a *App) User(r *http.Request) *UserInfo {
	c, err := r.Cookie(a.Cookie)
	if err != nil {
		return nil
	}

	u, err := FromJWT(c.Value, a.PublicKey, a.ID)
	if err != nil {
		return nil
	}
	return u
}

// A variant of http.Error
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err string, code int)

// Return a handler for the login. destination is the default destination.
// h is a handler on error case, if nil it replace by http.Error.
func (a *App) Login(destination string, h ErrorHandler) http.HandlerFunc {
	if destination == "" {
		destination = "/"
	}
	if h == nil {
		h = func(w http.ResponseWriter, r *http.Request, err string, code int) {
			http.Error(w, err, code)
		}
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
		if _, err := FromJWT(j, a.PublicKey, a.ID); err != nil {
			h(w, r, "JWT error: "+err.Error(), http.StatusBadRequest)
			return
		}

		c := model
		c.Value = j
		w.Header().Add("Set-Cookie", c.String())

		to := destination
		if s := fromb64(r.URL.Query().Get("r")); s != "" {
			to = s
		}
		redirection(w, to)
	}
}

// Return decoding base 64 Raw URL. If error, then output is empty.
func fromb64(in string) string {
	s, err := base64.RawURLEncoding.DecodeString(in)
	if err != nil {
		return ""
	}
	return string(s)
}

// Return the string encoding in base64
func tob64(in string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(in))
}

// Remove the cookie and redirect the client to the variable to; if it's
// empty, it's set to "/".
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

// Return the address top the provider to auth. r is teh address to go after
// autentification work.
func (a *App) Auth(r string) string {
	out := a.Server + "auth?app=" + a.ID
	if r != "" {
		return out + "&r=" + tob64(r)
	}
	return out
}

// Send JavaScrip to make redirection with cookies.
func redirection(w http.ResponseWriter, to string) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="utf-8"></head><body><a href="` + to + `">Redirect</a><script>document.location.replace('` + to + `');</script></body></html>`))
}
