// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package google

import (
	".."
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

var (
	// Redirect string
	ctx = context.Background()
	// Please config: ClientID, ClientSecret and RedirectURL
	Conf = oauth2.Config{
		Endpoint: google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}
)

func URL(to string) string {
	return "https://accounts.google.com/o/oauth2/v2/auth?" +
		"scope=https%3A//www.googleapis.com/auth/userinfo.profile%20https%3A//www.googleapis.com/auth/userinfo.email" +
		"&response_type=code" +
		"&client_id=" + Conf.ClientID +
		"&redirect_uri=" + Conf.RedirectURL
}

func User(r *http.Request) (*public.UserInfo, error) {
	u, err := getUser(r.URL.Query().Get("code"))
	if err != nil {
		return nil, err
	}

	return &public.UserInfo{
		ID:     "google:" + u.ID,
		Pseudo: u.Name,
		Email:  u.Email,
		Avatar: u.Picture,
	}, nil
}

type user struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// Download information from Google api
func getUser(code string) (*user, error) {
	t, err := Conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	client := Conf.Client(ctx, t)

	rep, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}

	var u user
	err = json.Unmarshal(all, &u)
	return &u, err
}
