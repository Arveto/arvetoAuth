// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"./pkg"
	"./pkg/github"
	"github.com/HuguesGuilleus/static.v1"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
)

func main() {
	static.Dev = true
	log.Println("main()")

	config, err := ini.Load("data/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	// Load github access
	github.Conf.ClientID = config.Section("github").Key("client").String()
	github.Conf.ClientSecret = config.Section("github").Key("secret").String()

	// Launch the server.
	log.Fatal(http.ListenAndServe(":8000", auth.Create(auth.Option{
		Key:          "data/key.pem",
		URL:          "http://localhost:8000/",
		MailHost:     config.Section("mail").Key("host").String(),
		MailLogin:    config.Section("mail").Key("login").String(),
		MailPassword: config.Section("mail").Key("password").String(),
	})))
}
