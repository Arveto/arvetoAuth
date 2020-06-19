// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

type Event struct {
	Id        string    `json:"id"`
	Actor     string    `json:"actor"`
	Operation string    `json:"operation"`
	Value     string    `json:"value"`
	Date      time.Time `json:"date"`
}

func (s *Server) logAdd(u *User, op, value string) {
	id := time.Now().Format("log:2006-01-02-05.000000")
	s.db.SetS(id, &Event{
		Actor:     u.Login,
		Operation: op,
		Value:     value,
		Id:        id,
		Date:      time.Now().UTC().Truncate(time.Minute),
	})
}

func (s *Server) logList(w http.ResponseWriter, r *http.Request) {
	// TODO: get year + mouth + day in prefix
	prefix := "log:"

	all := make([]Event, 0)
	s.db.ForS(prefix, 0, 0, nil, func(_ string, e Event) {
		all = append(all, e)
	})

	j, _ := json.Marshal(all)
	w.Header().Add("Content-Type", "application/json")
	w.Write(j)
}
