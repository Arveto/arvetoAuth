// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"./pkg"
	"./pkg/public/github"
	"./pkg/public/google"
	"github.com/HuguesGuilleus/go-logoutput"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
)

func init() {
	for _, a := range os.Args[1:] {
		if a == "-h" || a == "--help" {
			log.Println("(optionnal) Give the working directory")
			os.Exit(0)
		}
	}

	if len(os.Args) == 2 {
		err := os.Chdir(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	log.Println("[START]")

	config, err := ini.Load("data/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	url := config.Section("").Key("url").String()
	if l := config.Section("").Key("log").String(); l != "" {
		logoutput.SetLog(l)
	}

	// Load github access
	github.Conf.ClientID = config.Section("github").Key("client").String()
	github.Conf.ClientSecret = config.Section("github").Key("secret").String()
	google.Conf.ClientID = config.Section("google").Key("client").String()
	google.Conf.ClientSecret = config.Section("google").Key("secret").String()
	google.Conf.RedirectURL = url + "login/from/google/"

	// Launch the server.
	listen := config.Section("").Key("listen").MustString(":8000")
	log.Fatal(http.ListenAndServe(listen, auth.Create(auth.Option{
		Dev:          config.Section("").Key("dev").MustBool(),
		Key:          "data/key.pem",
		URL:          url,
		MailHost:     config.Section("mail").Key("host").String(),
		MailLogin:    config.Section("mail").Key("login").String(),
		MailPassword: config.Section("mail").Key("password").String(),
	})))
}
