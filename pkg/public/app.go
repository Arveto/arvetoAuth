package public

import (
	"crypto/rsa"
	"encoding/base64"
	"github.com/HuguesGuilleus/go-parsersa"
	"net/http"
)

// The application information. Fill all fields before use it.
type App struct {
	// The name of the app
	Name string
	// The public key of the auth server
	PublicKey *rsa.PublicKey
	// The name of the cookie who store teh JWT
	// by default it's `auth`
	Cookie string
}

// Create a new app. return an errro from open public key file.
func NewApp(name, pub, cookie string) (*App, error) {
	p, err := parsersa.PublicFile(pub)
	if err != nil {
		return nil, err
	}

	return &App{
		Name:      name,
		PublicKey: p,
		Cookie:    cookie,
	}, nil
}

// Get the user from a request cookie. If error occure, the function
// return nil. If you want the error details, use FromJWT inplace.
func (a *App) User(r *http.Request) *UserInfo {
	c, err := r.Cookie(a.Cookie)
	if err != nil {
		return nil
	}

	u, err := FromJWT(c.Value, a.PublicKey, a.Name)
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
	}

	return func(w http.ResponseWriter, r *http.Request) {
		j := r.URL.Query().Get("jwt")
		if _, err := FromJWT(j, a.PublicKey, a.Name); err != nil {
			h(w, r, "JWT error: "+err.Error(), http.StatusBadRequest)
			return
		}

		c := model
		c.Value = j
		c.Domain = r.Host
		w.Header().Add("Set-Cookie", c.String())

		to := destination
		if s := b64(r.URL.Query().Get("r")); s != "" {
			to = s
		}
		http.Redirect(w, r, to, http.StatusFound)
	}
}

// Return decoding base 64 Raw URL. If error, then output is empty.
func b64(in string) string {
	s, err := base64.RawStdEncoding.DecodeString(in)
	if err != nil {
		return ""
	}
	return string(s)
}

// Remove the cookie and redirect the client to the variable to; if it's
// empty, it's set to "/".
func (a *App) Logout(to string) http.HandlerFunc {
	if to == "" {
		to = "/"
	}

	c := (&http.Cookie{
		Name:   a.Cookie,
		Value:  "",
		MaxAge: -1,
	}).String()

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", c)
		http.Redirect(w, r, to, http.StatusFound)
	}
}
