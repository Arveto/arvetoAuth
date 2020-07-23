// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package public

import (
	"errors"
	"net/http"
)

// Manage request from login by a provider
type ExternProvider func(r *http.Request) (*UserInfo, error)

// Public informations, send to web page.
type UserInfo struct {
	ID     string    `json:"id"`
	Pseudo string    `json:"pseudo"`
	Email  string    `json:"email"`
	Avatar string    `json:"avatar"`
	Level  UserLevel `json:"level"`
}

type UserLevel int

const (
	LevelCandidate UserLevel = iota
	LevelVisitor   UserLevel = iota
	LevelStd       UserLevel = iota
	LevelAdmin     UserLevel = iota
	LevelBan       UserLevel = -1
)

func (ul UserLevel) String() string {
	switch ul {
	case LevelCandidate:
		return "Candidate"
	case LevelVisitor:
		return "Visitor"
	case LevelStd:
		return "Std"
	case LevelAdmin:
		return "Admin"
	case LevelBan:
		return "Ban"
	}
	return "?"
}

func (l UserLevel) MarshalText() ([]byte, error) {
	s := l.String()
	if s == "?" {
		return nil, errors.New("Unknown UserLevel")
	}
	return []byte(s), nil
}

func (l *UserLevel) UnmarshalText(text []byte) error {
	switch string(text) {
	case `Candidate`:
		*l = LevelCandidate
	case `Visitor`:
		*l = LevelVisitor
	case `Std`:
		*l = LevelStd
	case `Admin`:
		*l = LevelAdmin
	case `Ban`:
		*l = LevelBan
	default:
		return errors.New("Unknown UserLevel")
	}
	return nil
}
