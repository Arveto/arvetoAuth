// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Event struct {
	Id        string    `json:"id"`
	Actor     string    `json:"actor"`
	Operation string    `json:"operation"`
	Value     []string  `json:"value"`
	Date      time.Time `json:"date"`
}

// Add new event.
func (s *Server) logAdd(u *User, op string, value ...string) {
	id := time.Now().Format("log:2006-1-2-05.000000")
	s.db.SetS(id, &Event{
		Actor:     u.ID,
		Operation: op,
		Value:     value,
		Id:        id,
		Date:      time.Now().UTC().Truncate(time.Minute),
	})
}

// List the element for a specific period.
func (s *Server) logList(w http.ResponseWriter, r *http.Request) {
	prefix, err := logGetPrefix(w, r)
	if err {
		return
	}

	all := make([]Event, 0)
	s.db.ForS(prefix, 0, 0, nil, func(_ string, e Event) {
		all = append(all, e)
	})

	j, _ := json.Marshal(all)
	w.Header().Add("Content-Type", "application/json")
	w.Write(j)
}

// Send the number of element for a specifica period.
func (s *Server) logCount(w http.ResponseWriter, r *http.Request) {
	prefix, err := logGetPrefix(w, r)
	if err {
		return
	}

	n := 0
	s.db.ForS(prefix, 0, 0, func(string) bool {
		n++
		return false
	}, func(_ string, e Event) {})

	j, _ := json.Marshal(n)
	w.Header().Add("Content-Type", "application/json")
	w.Write(j)
}

// Get the prefix for log list or index operation
func logGetPrefix(w http.ResponseWriter, r *http.Request) (string, bool) {
	prefix := "log:"
	q := r.URL.Query()

	// Generic function to get a specific params
	getP := func(name, k string) int64 {
		if s := q.Get(k); s != "" {
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Read %s (params: %q) error:\r\n", name, k)
				fmt.Fprintln(w, err.Error())
				return -1
			}
			return v
		}
		return 0
	}

	// Year
	switch y := getP("year", "y"); y {
	case -1:
		return "", true
	case 0:
		prefix += time.Now().Format("2006-")
	default:
		prefix += strconv.FormatInt(y, 10) + "-"
	}

	// Mouth
	if m := getP("mouth", "m"); m < 0 {
		return "", true
	} else if m > 0 {
		prefix += strconv.FormatInt(m, 10) + "-"
		// Day
		if d := getP("day", "d"); d < 0 {
			return "", true
		} else if d > 0 {
			prefix += strconv.FormatInt(d, 10) + "-"
		}
	}

	return prefix, false
}
