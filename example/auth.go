// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"fmt"
	"github.com/Arveto/arvetoAuth/pkg/public"
	"net/http"
)

func main() {
	app, err := public.NewApp("ex", "...pub.pem", "jwt")
	if err != nil {
		fmt.Println("Load public key error:", err)
		return
	}

	// Standard page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[REQ]", r.URL)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		u := app.User(r)
		if u != nil {
			fmt.Fprintf(w, "%v\r\n", u)
			fmt.Fprintf(w, `<br><a href="/logout">Logout</a>`)
		} else {
			fmt.Fprintf(w, "You are not logged.<br>\r\n")
			fmt.Fprintf(w, `<a href="http://localhost:8000/auth?app=ex&r=%s">Login</a>`,
				base64.RawURLEncoding.EncodeToString([]byte(r.URL.String())))
		}
	})

	http.HandleFunc("/login", app.Login())
	http.HandleFunc("/logout", app.Logout(""))

	fmt.Println("http.ListenAndServe()")
	http.ListenAndServe(":9000", nil)
}
