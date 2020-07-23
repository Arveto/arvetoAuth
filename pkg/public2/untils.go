// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package public

import (
	"encoding/base64"
	"net/http"
	"strconv"
)

// UserInfo with some information to a specific service.
type jwtBody struct {
	UserInfo

	Audience       string `json:"aud"`
	ExpirationTime int64  `json:"exp"`
}

// Return decoding base 64 Raw URL. If error, then output is empty.
func fromb64(in string) string {
	s, err := base64.RawURLEncoding.DecodeString(in)
	if err != nil {
		return ""
	}
	return string(s)
}

// Return the string encoding in base64
func tob64(in []byte) string {
	return base64.RawURLEncoding.EncodeToString(in)
}

// Return the string encoding in base64
func tob64S(s string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

// Send JavaScrip to make redirection with cookies.
func redirection(w http.ResponseWriter, to string) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="utf-8"><script>document.location.replace(` + strconv.QuoteToASCII(to) + `);</script></head><body><a href="` + to + `">Redirect</a></body></html>`))
}
