// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"./public"
)

func (s *Server) loadDefaultUsers() {
	list := [...]public.UserInfo{
		{
			Pseudo: "Tomas Benedito Bonito",
			ID:     "tomas",
			Level:  public.LevelVisitor,
		},
		{
			Pseudo: "Zeck",
			ID:     "zeck",
			Level:  public.LevelStd,
		},
		{
			Pseudo: "Yuan Shikaï",
			ID:     "yuan",
			Level:  public.LevelStd,
		},
		{
			Pseudo: "Tzu Shikaï",
			ID:     "tzu",
			Level:  public.LevelAdmin,
		},
	}
	for _, u := range list {
		u.Email = u.ID + "@fai.mil"
		s.db.SetS("user:"+u.ID, User{
			UserInfo: u,
		})
	}
}
