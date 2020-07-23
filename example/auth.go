// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"fmt"
	"github.com/Arveto/arvetoAuth/pkg/public2"
	"net/http"
)

func main() {
	app, err := public.NewApp("app.example.com", "https://auth.example.com/", true)
	if err != nil {
		fmt.Println("Load public key error:", err)
		return
	}

	// Standard page
	app.HandleFunc("/", public.LevelCandidate, func(w http.ResponseWriter, r *public.Request) {
		fmt.Println("[REQ]", r.URL)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		if r.User != nil {
			fmt.Fprintf(w, "%+q<br>\r\n", r.User)
			fmt.Fprintf(w, "<img src='%s' alt='avatar'><br>\r\n", r.User.Avatar)
			fmt.Fprintf(w, `<a href="/logout">Logout</a>`)
		} else {
			fmt.Fprintf(w, "You are not logged.<br>\r\n")
			fmt.Fprintf(w, `<a href="%s">Login</a>`, app.ProviderAuth(r.RequestURI))
		}
	})

	fmt.Println("http.ListenAndServe()")
	http.ListenAndServe(":9000", nil)
}
