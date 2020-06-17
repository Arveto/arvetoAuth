// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package public

import (
	"errors"
)

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

func (ul UserLevel) MarshalJSON() ([]byte, error) {
	switch ul {
	case LevelCandidate:
		return []byte(`"Candidate"`), nil
	case LevelVisitor:
		return []byte(`"Visitor"`), nil
	case LevelStd:
		return []byte(`"Std"`), nil
	case LevelAdmin:
		return []byte(`"Admin"`), nil
	case LevelBan:
		return []byte(`"Ban"`), nil
	}
	return nil, errors.New("Unknown UserLevel")
}

func (ul *UserLevel) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"Candidate"`:
		*ul = LevelCandidate
	case `"Visitor"`:
		*ul = LevelVisitor
	case `"Std"`:
		*ul = LevelStd
	case `"Admin"`:
		*ul = LevelAdmin
	case `"Ban"`:
		*ul = LevelBan
	default:
		return errors.New("Unknown UserLevel")
	}
	return nil
}
