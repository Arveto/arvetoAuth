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
			Name:  "Tomas Benedito Bonito",
			Login: "tomas",
			Level: public.LevelStd,
		},
		{
			Name:  "Zeck",
			Login: "zeck",
			Level: public.LevelStd,
		},
		{
			Name:  "Yuan Shikaï",
			Login: "yuan",
			Level: public.LevelStd,
		},
		{
			Name:  "Tzu Shikaï",
			Login: "tzu",
			Level: public.LevelAdmin,
		},
	}
	for _, u := range list {
		u.Email = u.Login + "@fai.mil"
		s.db.SetS("user:"+u.Login, User{
			UserInfo: u,
		})
	}
}
