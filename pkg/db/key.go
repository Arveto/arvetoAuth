// Copyright (c) 2020, HuguesGuilleus. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"log"
	"net/http"
)

// The number index of an element in a DB.
type Key uint32

/* STRING */

func (k Key) String() string {
	return fmt.Sprintf("%d", uint(k))
}

func KeyFromString(s string) (k Key) {
	if s == "" {
		return Key(0)
	}
	_, err := fmt.Sscanf(s, "%d", &k)
	if err != nil {
		log.Println(err)
	}
	return
}

// Get the key from the parametre q of the request URL
func KeyFromReq(r *http.Request, q string) Key {
	return KeyFromString(r.URL.Query().Get(q))
}

/* BYTES */

func (k Key) bytes() (b []byte) {
	if k == 0 {
		return []byte{0}
	}
	for k != 0 {
		b = append(b, byte(k&0xFF))
		k /= 0x100
	}
	return
}

func keyBytes(bytes []byte) (k Key) {
	if l := len(bytes); l > 4 {
		bytes = bytes[l-4:]
	}
	for _, b := range bytes {
		k = k*0x100 + Key(b)
	}
	return
}
