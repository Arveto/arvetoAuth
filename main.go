// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"./pkg"
	"./pkg/github"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
)

// Load github access
func init() {
	config, err := ini.Load("data/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	github.Client = config.Section("github").Key("client").String()
	github.Secret = config.Section("github").Key("secret").String()
}

func main() {
	log.Println("main()")
	log.Fatal(http.ListenAndServe(":8000", auth.Create(auth.Option{
		Key:  "data/key.pem",
		URL: "http://localhost:8000/",
	})))
}
