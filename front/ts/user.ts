// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace User {
	interface User {
		login: string;
		pseudo: string;
		email: string;
		avatar: string;
		level: string;
	}

	export async function list() {
		let t = Deskop.table(['', 'ID', 'Pseudo', 'Email', 'Level']);
		let l: User[] = await (await fetch('/user/list')).json();
		console.log("l:", l);

		l.forEach(u => t.insertAdjacentHTML('beforeend', `<tr>
			<td></td>
			<td>${u.login}</td>
			<td>${u.pseudo}</td>
			<td>${u.email}</td>
			<td>${u.level}</td>
		</tr>`))
	}
}
