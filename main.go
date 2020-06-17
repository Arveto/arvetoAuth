// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"./pkg"
	"log"
	"net/http"
)

func main() {
	log.Println("main()")
	log.Fatal(http.ListenAndServe(":8000", auth.Create(auth.Option{
		Key: "data/key.pem",
	})))
}
