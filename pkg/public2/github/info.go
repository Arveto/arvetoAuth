// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package github

import (
	".."
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io/ioutil"
	"net/http"
)

var (
	loginURL = ""
	ctx      = context.Background()
	// Please config: ClientID and ClientSecret
	Conf = oauth2.Config{
		Endpoint: github.Endpoint,
	}
)

// Return and URL to come to GitHub signing page.
func URL(redirect string) string {
	if loginURL == "" {
		loginURL = "https://github.com/login/oauth/authorize?" +
			"scope=user%3Aemail" +
			"&client_id=" + Conf.ClientID
	}

	if redirect != "" {
		return loginURL + "&redirect_uri=" + redirect
	}

	return loginURL
}

type Info struct {
	Pseudo string `json:"name"`
	Login  string `json:"login"`
	Email  string `json:"email"`
	Avatar string `json:"avatar_url"`
}

// Return an user login with GitHub
func User(r *http.Request) (*public.UserInfo, error) {
	i, err := downloadInfo(r)
	if err != nil {
		return nil, err
	}

	return &public.UserInfo{
		ID:     "github:" + i.Login,
		Pseudo: i.Pseudo,
		Email:  i.Email,
		Avatar: i.Avatar,
	}, nil
}

func downloadInfo(r *http.Request) (*Info, error) {
	t, err := Conf.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		return nil, err
	}
	client := Conf.Client(ctx, t)

	// Get information
	i, err := getInfo(client)
	if err != nil {
		return nil, err
	}
	// Get email adress
	i.Email, err = getMail(client)
	if i.Pseudo == "" {
		i.Pseudo = i.Login
	}

	return i, err
}

func getInfo(client *http.Client) (*Info, error) {
	resp, err := client.Get("https://api.github.com/user")
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
