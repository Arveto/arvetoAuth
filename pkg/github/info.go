package github

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	Client   = ""
	Secret   = ""
	loginURL = ""
)

// Return and URL to come to GitHub signing page.
func URL(redirect string) string {
	if loginURL == "" && Client != "" {
		loginURL = "https://github.com/login/oauth/authorize?" +
			"scope=user%3Aemail" +
			"&client_id=" + Client
	}

	if redirect != "" {
		return loginURL + "&redirect_uri=" + redirect
	}

	return loginURL
}

type Info struct {
	Pseudo   string `json:"name"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Icon     string `json:"avatar_url"`
	GithubId int    `json:"id"`
}

func NewInfo(r *http.Request) (*Info, error) {
	// Get tocken
	token, err := loginGetToken(r.URL.Query().Get("code"))
	if err != nil {
		return nil, err
	}

	// Get information
	i, err := getInfo(token)
	if err != nil {
		return nil, err
	}
	i.Email, err = getMail(token)
	if i.Pseudo == "" {
		i.Pseudo = i.Login
	}

	return i, err
}

func getInfo(token string) (*Info, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var i Info
	err = json.Unmarshal(data, &i)
	return &i, err
}

// Get the token with user code
func loginGetToken(code string) (string, error) {
	q := url.Values{}
	q.Add("client_id", Client)
	q.Add("client_secret", Secret)
	q.Add("code", code)

	rep, err := http.PostForm("https://github.com/login/oauth/access_token", q)
	if err != nil {
		return "", err
	}
	defer rep.Body.Close()

	all, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return "", err
	}

	values, err := url.ParseQuery(string(all))
	if err != nil {
		return "", err
	}

	return values.Get("access_token"), nil
}
